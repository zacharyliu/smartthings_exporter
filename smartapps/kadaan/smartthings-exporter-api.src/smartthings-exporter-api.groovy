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

definition(
  name: "Smartthings_exporter API",
  namespace: "kadaan",
  author: "Joel Baranick",
  description: "API used by Smartthings_exporter to read sensor data.",
  category: "My Apps",
  iconUrl: "https://s3.amazonaws.com/smartapp-icons/Convenience/Cat-Convenience.png",
  iconX2Url: "https://s3.amazonaws.com/smartapp-icons/Convenience/Cat-Convenience@2x.png",
  oauth: [displayName: "Smartthings_exporter API", displayLink: ""])

preferences {
  section() {
    paragraph "Select the sensors you want the API to have access to."
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
  def result = [:]
  def enumValues = [:]
  def sensorResults = []
  sensorResults << sensors.collect{getSensor(enumValues, it)}
  result << ["enumValues": enumValues]
  result << ["sensors": sensorResults[0]]
  result
}

private getMapping(string_ordinal_map, ordinal_string_map, value) {
  if (value instanceof String) {
    def ordinal = ordinal_string_map.size()
    if (!string_ordinal_map.containsKey(value)) {
      ordinal_string_map << [(ordinal) : value]
      string_ordinal_map << [(value) : ordinal]
    }
    return ordinal
  } else {
    return value
  }
}

private isAllowed(attributeName) {
  if (attributeName == "trackData") {
    return false
  }
  return true
}

private getSensor(enumValues, sensor) {
  if (!sensor) return null
  def results = [:]
  ["id", "name", "displayName"].each {
    results << [(it) : sensor."$it"]
  }

  def attrsAndVals = [:]
  sensor.supportedAttributes?.each {
    if (isAllowed(it.name)) { 
      if (it.getDataType() == "ENUM" && !enumValues.containsKey(it.name)) {
        enumValues << [(it.name) : it.getValues()]
      }
      attrsAndVals << [(it.name) : sensor.currentValue(it.name)]
    }
  }
  results << ["attributes" : attrsAndVals]
  results
}