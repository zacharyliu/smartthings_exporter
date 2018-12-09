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

import groovy.transform.EqualsAndHashCode

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
  def strings = [:]
  def map = []
  def sensorResults = []
  sensorResults << getSensors(strings, map, sensors)
  result << ["m": map]
  result << ["s": sensorResults[0]]
  result
}

private isAllowedSensor(sensorType) {
  if (sensorType == "LAN Sonos Player") {
    return false
  }
  return true
}

private getMappedString(strings, map, value) {
  if (strings.containsKey(value)) {
    return strings[value]
  }
  def ordinal = strings.size()
  strings[value] = ordinal
  map << value
  return ordinal
}

private getSensors(strings, map, sensors) {
  def groupedSensors = sensors.groupBy({ [typeName:it.typeName, capabilities:it.capabilities.join(",")] })
  groupedSensors.findResults{k, v -> getSensor(strings, map, v)}
}

private getSensor(strings, map, sensors) {
  def sensor = sensors[0]
  if (!isAllowedSensor(sensor.getTypeName())) return null
  def results = []
  results << getMappedString(strings, map, sensor.getName())
  def size = sensors.size()
  if (size == 1) {
    def mappedValue = getMappedString(strings, map, sensor.getDisplayName())
    mappedValue = -mappedValue
    results << mappedValue
  } else {
    results << sensors.size()
    sensors.each {sens ->
      results << getMappedString(strings, map, sens.getDisplayName())
    }
  }
  sensor.capabilities?.each {cap ->
    cap.attributes?.each {attr ->
      if (sensors.any{s -> s.currentValue(attr.name) != null}) {
        def mappedAttr = getMappedString(strings, map, attr.name)
        switch (attr.getDataType()) {
          case "NUMBER":
            results << mappedAttr
			sensors.each {s -> 
              results << s.currentValue(attr.name)
            }
            break
          case "ENUM":
          	def written_enum_values = [:]
            mappedAttr = -mappedAttr
            results << mappedAttr
            def enumValues = attr.getValues()
            results << enumValues.size()
            sensors.each {s -> 
              def value = s.currentValue(attr.name)
		      enumValues.each {enumVal ->
                def mappedValue = getMappedString(strings, map, enumVal)
                if (value == enumVal) {
                  written_enum_values << [(mappedValue): true]
                  mappedValue = -mappedValue
                  results << mappedValue
                } else if (!written_enum_values.containsKey(mappedValue)) {
                  written_enum_values << [(mappedValue): true]
                  results << mappedValue
                }
              }
            }
            break
        }
      }
    }
  }
  results
}