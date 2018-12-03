// Copyright © 2018 Joel Baranick <jbaranick@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Based on:
// http://github.com/marcopaganini/smartcollector
// (C) 2016 by Marco Paganini <paganini@paganini.net>

package main

import (
	"fmt"
	"github.com/kadaan/gosmart"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	plog "github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"path/filepath"
)

const (
	namespace       = "smartthings"
	tokenFilePrefix = ".smartthings"
)

var (
	alarmState        *prometheus.Desc
	battery           *prometheus.Desc
	carbonMonoxide    *prometheus.Desc
	contact           *prometheus.Desc
	energy            *prometheus.Desc
	motion            *prometheus.Desc
	power             *prometheus.Desc
	presence          *prometheus.Desc
	smoke             *prometheus.Desc
	smartSwitch       *prometheus.Desc
	temperature       *prometheus.Desc
	valOpenClosed     = []string{"open", "closed"}
	valInactiveActive = []string{"inactive", "active"}
	valAbsentPresent  = []string{"not present", "present"}
	valOffOn          = []string{"off", "on"}
	mappingFunctions  = map[string]func(interface{}) (*prometheus.Desc, float64, error){
		"alarmState":     getAlarmState,
		"battery":        getBattery,
		"carbonMonoxide": getCarbonMonoxide,
		"contact":        getContact,
		"energy":         getEnergy,
		"motion":         getMotion,
		"power":          getPower,
		"presence":       getPresence,
		"smoke":          getSmoke,
		"switch":         getSwitch,
		"temperature":    getTemperature,
	}
)

// Exporter collects smartthings stats and exports them using the prometheus metrics package.
type Exporter struct {
	client   *http.Client
	endpoint string
}

// NewExporter returns an initialized Exporter.
func NewExporter(oauthClient string, oauthSecret string) (*Exporter, error) {
	tfile := tokenFilename(oauthClient)

	// Create the oauth2.config object and get a token
	config := gosmart.NewOAuthConfig(oauthClient, oauthSecret)
	token, err := gosmart.GetToken(tfile, config)
	if err != nil {
		plog.Fatalf("Error fetching token: %v", err)
	}

	// Create a client with the token and fetch endpoints URI.
	ctx := context.Background()
	client := config.Client(ctx, token)
	endpoint, err := gosmart.GetEndPointsURI(client)
	if err != nil {
		plog.Fatalf("Error reading endpoints URI: %v\n", err)
	}

	_, verr := gosmart.GetDevices(client, endpoint)
	if verr != nil {
		plog.Fatalf("Error verifying connection to endpoints URI: %v\n", err)
	}

	// Init our exporter.
	return &Exporter{
		client:   client,
		endpoint: endpoint,
	}, nil
}

// Describe describes all the metrics ever exported by the Kafka exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- alarmState
	ch <- battery
	ch <- carbonMonoxide
	ch <- contact
	ch <- energy
	ch <- motion
	ch <- power
	ch <- presence
	ch <- smoke
	ch <- smartSwitch
	ch <- temperature
}

// Collect fetches the stats from configured Kafka location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Iterate over all devices and collect timeseries info.
	devs, err := gosmart.GetDevices(e.client, e.endpoint)
	if err != nil {
		plog.Fatalf("Error reading list of devices: %v\n", err)
	}

	for _, dev := range devs {
		for k, val := range dev.Attributes {
			if val == nil {
				val = ""
			}
			var value float64
			var metricDesc *prometheus.Desc
			mappingFunction := mappingFunctions[k]
			if mappingFunction == nil {
				continue
			}
			metricDesc, value, err = mappingFunction(val)
			if err == nil {
				ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, value, dev.ID, dev.DisplayName)
			} else {
				plog.Errorf("Cannot process sensor data for %s: %v", k, err)
			}
		}
	}
}

func getAlarmState(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueFloat(val)
	return battery, value, err
}

func getBattery(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueClear(val)
	return alarmState, value, err
}

func getCarbonMonoxide(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueClear(val)
	return carbonMonoxide, value, err
}

func getContact(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueOneOf(val, valOpenClosed)
	return contact, value, err
}

func getEnergy(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueFloat(val)
	return energy, value, err
}

func getMotion(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueOneOf(val, valInactiveActive)
	return motion, value, err
}

func getPower(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueFloat(val)
	return power, value, err
}

func getPresence(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueOneOf(val, valAbsentPresent)
	return presence, value, err
}

func getSmoke(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueClear(val)
	return smoke, value, err
}

func getSwitch(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueOneOf(val, valOffOn)
	return smartSwitch, value, err
}

func getTemperature(val interface{}) (*prometheus.Desc, float64, error) {
	value, err := valueFloat(val)
	return temperature, value, err
}

// valueClear expects a string and returns 0 for "clear", 1 for anything else.
// TODO: Expand this to properly identify non-clear conditions and return error
// in case an unexpected value is found.
func valueClear(v interface{}) (float64, error) {
	val, ok := v.(string)
	if !ok {
		return 0.0, fmt.Errorf("invalid non-string argument %v", v)
	}
	if val == "clear" {
		return 0.0, nil
	}
	return 1.0, nil
}

// valueOneOf returns 0.0 if the value matches the first item
// in the array, 1.0 if it matches the second, and an error if
// nothing matches.
func valueOneOf(v interface{}, options []string) (float64, error) {
	val, ok := v.(string)
	if !ok {
		return 0.0, fmt.Errorf("invalid non-string argument %v", v)
	}
	if val == options[0] {
		return 0.0, nil
	}
	if val == options[1] {
		return 1.0, nil
	}
	return 0.0, fmt.Errorf("invalid option %q. Expected %q or %q", val, options[0], options[1])
}

// valueFloat returns the float64 value of the value passed or
// error if the value cannot be converted.
func valueFloat(v interface{}) (float64, error) {
	val, ok := v.(float64)
	if !ok {
		return 0.0, fmt.Errorf("invalid non floating-point argument %v", v)
	}
	return val, nil
}

func tokenFilename(oauthClient string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	return filepath.FromSlash(exPath + "/" + tokenFilePrefix + "_" + oauthClient + ".json")
}

func init() {
	prometheus.MustRegister(version.NewCollector("smartthings_exporter"))
}

func needsSecret(oauthClient string) bool {
	tfile := tokenFilename(oauthClient)
	token, err := gosmart.LoadToken(tfile)
	if err != nil {
		return true
	}
	return !token.Valid()
}

func main() {
	application := kingpin.New("smartthings_exporter", "Smartthings exporter for Prometheus")
	smartthingsOAuthSecretFlag := application.Flag("smartthings.oauth-secret", "Smartthings OAuth secret key.").Default("")
	var (
		listenAddress          = application.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9499").String()
		metricsPath            = application.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		smartthingsOAuthClient = application.Flag("smartthings.oauth-client", "Smartthings OAuth client ID.").Required().String()
		smartthingsOAuthSecret = smartthingsOAuthSecretFlag.String()
	)

	plog.AddFlags(application)
	application.Validate(func(application *kingpin.Application) error {
		if needsSecret(*smartthingsOAuthClient) && *smartthingsOAuthSecret == "" {
			return fmt.Errorf("required flag --%s not provided", smartthingsOAuthSecretFlag.Model().Name)
		}
		return nil
	})
	application.Version(version.Print("smartthings_exporter"))
	application.HelpFlag.Short('h')
	_, err := application.Parse(os.Args[1:])
	if err != nil {
		application.Fatalf("%s, try --help", err)
	}

	plog.Infoln("Starting smartthings_exporter", version.Info())
	plog.Infoln("Build context", version.BuildContext())

	alarmState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "alarm_clear"),
		"0 if the alarm is clear.",
		[]string{"id", "name"}, nil,
	)
	battery = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "battery_percentage"),
		"Percentage of battery remaining.",
		[]string{"id", "name"}, nil,
	)
	carbonMonoxide = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "carbonMonoxide_clear"),
		"0 if no carbon monoxide is detected.",
		[]string{"id", "name"}, nil,
	)
	contact = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "contact_closed"),
		"1 if the contact is closed.",
		[]string{"id", "name"}, nil,
	)

	energy = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "energy_usage_kilowatts"),
		"Energy usage in kilowatts.",
		[]string{"id", "name"}, nil,
	)

	motion = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "motion_detected"),
		"1 if presence is detected.",
		[]string{"id", "name"}, nil,
	)

	power = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "power_usage_watts"),
		"Current power usage in watts.",
		[]string{"id", "name"}, nil,
	)

	presence = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "presence_detected"),
		"1 if presence is detected.",
		[]string{"id", "name"}, nil,
	)

	smoke = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "smoke_detected"),
		"1 if smoke is detected.",
		[]string{"id", "name"}, nil,
	)

	smartSwitch = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "switch_on"),
		"1 if the switch is on.",
		[]string{"id", "name"}, nil,
	)

	temperature = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "temperature_f"),
		"Temperature in fahrenheit.",
		[]string{"id", "name"}, nil,
	)

	exporter, err := NewExporter(*smartthingsOAuthClient, *smartthingsOAuthSecret)
	if err != nil {
		plog.Fatalln(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
	        <head><title>SmartThings Exporter</title></head>
	        <body>
	        <h1>SmartThings Exporter</h1>
	        <p><a href='` + *metricsPath + `'>Metrics</a></p>
	        </body>
	        </html>`))
	})

	plog.Infoln("Listening on", *listenAddress)
	plog.Fatal(http.ListenAndServe(*listenAddress, nil))
}
