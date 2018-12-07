/**
 *  Copyright 2015 SmartThings
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at:
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
 *  on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License
 *  for the specific language governing permissions and limitations under the License.
 *
 *  Developer API
 *
 *  Author: SmartThings
 */

import groovy.json.JsonBuilder

definition(
  name: "Sensor API",
  namespace: "kadaan",
  author: "Joel Baranick",
  description: "Sensor API used by Smartthings_exporter to read sensor data.",
  category: "My Apps",
  iconUrl: "https://s3.amazonaws.com/smartapp-icons/Convenience/Cat-Convenience.png",
  iconX2Url: "https://s3.amazonaws.com/smartapp-icons/Convenience/Cat-Convenience@2x.png",
  oauth: [displayName: "Device API", displayLink: ""])

preferences {
  section() {
    paragraph "Select the sensors you want API to have access to."
  }
  section() {
  	input "sensors", "capability.sensor", multiple: true, title: "Which sensors?", required: true
  }
}

def installed() {
  initialize()
}

def updated() {
  unsubscribe()
  initialize()
}

def initialize() {
}

mappings {
  path("/sensors") {
    action: [
      GET: "listSensors"
    ]
  }
}

def listSensors() {
  def result = []
  result << sensors.collect{getSensor(it)}
  log.debug "Returning DEVICES: $result"
  result[0]
}

private getSensor(sensor) {
  if (!sensor) return null
  def results = [:]
  ["id", "name", "displayName"].each {
    results << [(it) : sensor."$it"]
  }

  def attrsAndVals = [:]
  sensor.supportedAttributes?.each {
    attrsAndVals << [(it.name) : sensor.currentValue(it.name)]
  }

  results << ["attributes" : attrsAndVals]
  log.debug "Returning DEVICE: $results"
  results
}