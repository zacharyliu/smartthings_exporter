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
  name: "${handle()} API",
  namespace: "kadaan",
  author: "Joel Baranick",
  description: "Tap here to install ${handle()} ${version()}",
  parent: "kadaan:${handle()}",
  category: "My Apps",
  iconUrl: "https://s3.amazonaws.com/smartapp-icons/Convenience/Cat-Convenience.png",
  iconX2Url: "https://s3.amazonaws.com/smartapp-icons/Convenience/Cat-Convenience@2x.png",
  oauth: [displayName: "smartthings_exporter API", displayLink: ""])

preferences {
  page(name: "pageSettings")
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
    path("/devices") {
        action: [
            GET: "listDevices"
        ]
    }
}

def listDevices() {
    result << allDevices.collect{deviceItem(it)}
    log.debug "Returning DEVICES: $result"
    result[0]
}

private deviceItem(device) {
    if (!device) return null
    def results = [:]
    ["id", "name", "displayName"].each {
        results << [(it) : device."$it"]
    }

    def attrsAndVals = [:]
    device.supportedAttributes?.each {
        attrsAndVals << [(it.name) : device.currentValue(it.name)]
    }

    results << ["attributes" : attrsAndVals]
    log.debug "Returning DEVICE: $results"
    results
}

def pageSettings() {
    //clear devices cache
  dynamicPage(name: "pageSettings", title: "", install: false, uninstall: false) {
    section("Available devices") {
      href "pageSelectDevices", title: "Available devices", description: "Tap here to select which devices are available to smartthings_exporter"
    }
  }
}

private pageSelectDevices() {
  state.deviceVersion = now().toString()
    dynamicPage(name: "pageSelectDevices", title: "") {
    section() {
      paragraph "Select the devices you want smartthings_exporter to have access to."
    }
    section ('Select devices by type') {
      paragraph "Most devices should fall into one of these two categories"
      input "dev:actuator", "capability.actuator", multiple: true, title: "Which actuators", required: false
      input "dev:sensor", "capability.sensor", multiple: true, title: "Which sensors", required: false
    }

    section ('Select devices by capability') {
      paragraph "If you cannot find a device by type, you may try looking for it by category below"
      def d
      for (capability in capabilities().findAll{ (!(it.value.d in [null, 'actuators', 'sensors'])) }.sort{ it.value.d }) {
        if (capability.value.d != d) input "dev:${capability.key}", "capability.${capability.key}", multiple: true, title: "Which ${capability.value.d}", required: false
        d = capability.value.d
      }
    }
  }
}