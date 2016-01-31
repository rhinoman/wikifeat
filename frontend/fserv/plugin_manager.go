/*
 *  Licensed to Wikifeat under one or more contributor license agreements.
 *  See the LICENSE.txt file distributed with this work for additional information
 *  regarding copyright ownership.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *  notice, this list of conditions and the following disclaimer in the
 *  documentation and/or other materials provided with the distribution.
 *  * Neither the name of Wikifeat nor the names of its contributors may be used
 *  to endorse or promote products derived from this software without
 *  specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

package fserv

import (
	"errors"
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/alyu/configparser"
	"log"
)

type PluginData struct {
	Name       string `json:"id"`
	Author     string `json:"author"`
	Version    string `json:"version"`
	PluginDir  string `json:"pluginDir"`
	MainScript string `json:"mainScript"`
	Stylesheet string `json:"stylesheet"`
	Enabled    bool   `json:"enabled"`
}

var enabledPlugins = make(map[string]PluginData)

// Return the Enabled Plugins map
func GetEnabledPlugins() []PluginData {
	pluginList := []PluginData{}
	for _, v := range enabledPlugins {
		pluginList = append(pluginList, v)
	}
	return pluginList
}

// Fetches an item from the plugins map
func GetPluginData(pluginName string) (*PluginData, error) {
	if pg, ok := enabledPlugins[pluginName]; ok == false {
		return nil, errors.New("Not Found")
	} else {
		return &pg, nil
	}
}

// Creates a new Default PluginData struct
func NewPluginData() *PluginData {
	return &PluginData{
		Name:       "Unnamed",
		Author:     "Unknown",
		Version:    "0.1",
		PluginDir:  "unnamed",
		MainScript: "main.js",
	}
}

// Loads Plugin Data from the plugins ini file
func LoadPluginData(filename string) {
	readSinglePluginData := func(pluginSection *configparser.Section) *PluginData {
		theData := NewPluginData()
		for key, value := range pluginSection.Options() {
			switch key {
			case "name":
				theData.Name = value
			case "author":
				theData.Author = value
			case "version":
				theData.Version = value
			case "pluginDir":
				theData.PluginDir = value
			case "mainScript":
				theData.MainScript = value
			case "stylesheet":
				theData.Stylesheet = value
			case "enabled":
				if value == "true" {
					theData.Enabled = true
				} else {
					theData.Enabled = false
				}
			}
		}
		return theData
	}
	config, err := configparser.Read(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(config)
	pluginSections, err := config.Find(".plugin$")
	if err != nil {
		log.Fatal(err)
	}
	for _, section := range pluginSections {
		theData := readSinglePluginData(section)
		if theData.Enabled {
			enabledPlugins[theData.Name] = *theData
		}
	}
}
