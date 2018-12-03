# smarthings_exporter [![Build Status](https://travis-ci.org/kadaan/smarthings_exporter.svg?branch=master)](https://travis-ci.org/kadaan/smarthings_exporter) [![Coverage Status](https://img.shields.io/coveralls/github/kadaan/smarthings_exporter/master.svg)](https://coveralls.io/github/kadaan/smarthings_exporter) [![Go Report Card](https://goreportcard.com/badge/github.com/kadaan/smarthings_exporter)](https://goreportcard.com/report/github.com/kadaan/smarthings_exporter)

smartthings_exporter is a command line tool to export information about your SmartThings
sensors in a format that can be scraped by [Prometheus](http://prometheus.io). The tool uses the [GoSmart](http://github.com/marcopaganini/gosmart) library to talk to SmartThings and collect sensor data and exposed the metrics over http.

## Installation

The installation instructions assume a properly installed and configured Go
development environment. The very first step is to download and build
smartthings_exporter (this step will also download and compile the GoSmart library):


```
$ go get -u github.com/kadaan/smartthings_exporter
```

### GoSmart configuration

Before we can use smarthings_exporter, we need to configure GoSmart. Follow the
[GoSmart installation instructions](https://github.com/kadaan/gosmart#installation)
carefully, making sure all steps have been followed.

With GoSmart configured,
[Follow the instructions](https://github.com/kadaan/gosmart#running-an-example) to
run the simple example that comes with GoSmart. Make sure the example displays
a list of your sensors on the screen.

### smartthings_exporter configuration

We now need to authorize smartthings_exporter to access your Smartthings app. Take note of the `client_id` and `client_secret` of your SmartThings app (used when running the simple example above). Run:

```
$ smartthings_exporter --smartthings.oauth-client=<client> --smartthings.oauth-secret=<secret>
```

Follow the prompts to authorize the app.

smartthings_exporter will write a file with your credentials to executable directory. After the first run, only the `--smartthings.oauth-client` is required to run smarthings_exporter:

```
$ smarthings_exporter --smartthings.oauth-client=<client_id>
```