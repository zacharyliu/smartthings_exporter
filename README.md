# smarthings_exporter [![Build Status](https://travis-ci.org/kadaan/smartthings_exporter.svg?branch=master)](https://travis-ci.org/kadaan/smartthings_exporter) [![Coverage Status](https://img.shields.io/coveralls/github/kadaan/smartthings_exporter/master.svg)](https://coveralls.io/github/kadaan/smartthings_exporter) [![Go Report Card](https://goreportcard.com/badge/github.com/kadaan/smartthings_exporter)](https://goreportcard.com/report/github.com/kadaan/smartthings_exporter)

Smartthings_exporter is a command line tool to export information about your SmartThings
sensors in a format that can be scraped by [Prometheus](http://prometheus.io). The tool uses 
the [GoSmart](http://github.com/marcopaganini/gosmart) library to talk to SmartThings and collect 
sensor data and exposed the metrics over http.

## Installation

The installation instructions assume a properly installed and configured Go
development environment. The very first step is to download and build
Smartthings_exporter (this step will also download and compile the GoSmart library):


```
$ go get -u github.com/kadaan/smartthings_exporter
```

### Smartthings Setup

Before you can use Smartthings_exporter, you need to register it with Smartthings.  

The first step is to setup the API that Smartthings_exporter uses to communicate with Smartthings.  Follow the 
[GoSmart Smartthings API setup](https://github.com/kadaan/gosmart#smartthings-api-setup) steps.

Take note of the `client_id` and `client_secret` of your SmartThings app that you just created.

### Smartthings_exporter configuration

We now need to register Smartthings_exporter to with your Smartthings app.

Run:

```
$ smartthings_exporter register --smartthings.oauth-client=[client_id] --smartthings.oauth-secret==[client_secret] > .st_token
```

Follow the prompts to authorize the app.

## Running

Now we can start Smartthings_exporter by running:

```
$ smartthings_exporter --smartthings.oauth-client=<client_id> --smartthings.oauth-token.file=.sttoken
```

# Prometheus Operator

```bash
kubectl port-forward service/prometheus-operator-prometheus 9090
kubectl port-forward service/prometheus-operator-grafana 8080:80
```

Useful links:
- https://coreos.com/blog/the-prometheus-operator.html
- https://medium.com/faun/trying-prometheus-operator-with-helm-minikube-b617a2dccfa3
- https://grafana.com/docs/features/datasources/prometheus/
- https://blog.kubernauts.io/cloud-native-monitoring-with-prometheus-and-grafana-9c8003ab9c7
