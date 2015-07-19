/**  Copyright (c) 2014-present James Adam.  All rights reserved.
*
*		 This file is part of Wikifeat.
*
*    Wikifeat is free software: you can redistribute it and/or modify
*    it under the terms of the GNU General Public License as published by
*    the Free Software Foundation, either version 2 of the License, or
*    (at your option) any later version.
*
*    This program is distributed in the hope that it will be useful,
*    but WITHOUT ANY WARRANTY; without even the implied warranty of
*    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*    GNU General Public License for more details.
*
*    You should have received a copy of the GNU General Public License
*    along with this program.  If not, see <http://www.gnu.org/licenses/>.
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

// Creates a new Default PLuginData struct
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
