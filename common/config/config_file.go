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

// Stores configuration for the core services
package config

import (
	"fmt"
	"github.com/alyu/configparser"
	"github.com/rhinoman/wikifeat/common/util"
	"log"
	"path"
	"strconv"
	"strings"
)

// Load config values from file
func LoadConfig(filename string) {
	LoadDefaults()
	log.Printf("\nLoading Configuration from %v\n", filename)
	config, err := configparser.Read(filename)
	fmt.Print(config)
	if err != nil {
		log.Fatal(err)
	}
	dbSection, err := config.Section("Database")
	if err != nil {
		log.Fatal(err)
	}
	logSection, err := config.Section("Logging")
	if err != nil {
		log.Fatal(err)
	}
	authSection, err := config.Section("Auth")
	if err != nil {
		log.Fatal(err)
	}
	registrySection, err := config.Section("ServiceRegistry")
	if err != nil {
		log.Fatal(err)
	}
	//Optional sections
	frontendSection, err := config.Section("Frontend")
	userSection, err := config.Section("Users")
	searchSection, err := config.Section("Search")
	notifSection, err := config.Section("Notifications")
	if frontendSection != nil {
		SetFrontendConfig(frontendSection)
	}
	if searchSection != nil {
		SetSearchConfig(searchSection)
	}
	if userSection != nil {
		setUsersConfig(userSection)
	}
	if notifSection != nil {
		setNotificationConfig(notifSection)
	}
	setDbConfig(dbSection)
	setLogConfig(logSection)
	setAuthConfig(authSection)
	setRegistryConfig(registrySection)
}

// Load Frontend configuration options
func SetFrontendConfig(frontendSection *configparser.Section) {
	execDir, _ := util.GetExecDirectory()
	for key, value := range frontendSection.Options() {
		switch key {
		case "webAppDir":
			if value[0] != '/' {
				Frontend.WebAppDir = path.Join(execDir, value)
			} else {
				Frontend.WebAppDir = value
			}
		case "pluginDir":
			if value[0] != '/' {
				Frontend.PluginDir = path.Join(execDir, value)
			} else {
				Frontend.PluginDir = value
			}
		case "homepage":
			Frontend.Homepage = value
		}
	}
}

// Load Search configuration options
func SetSearchConfig(searchSection *configparser.Section) {
	for key, value := range searchSection.Options() {
		switch key {
		case "searchServerLocation":
			Search.SearchServerLocation = value
		}
	}
}

// Load Database configuration options
func setDbConfig(dbSection *configparser.Section) {
	for key, value := range dbSection.Options() {
		switch key {
		case "dbAddr":
			Database.DbAddr = value
		case "dbPort":
			Database.DbPort = value
		case "useSSL":
			Database.UseSSL = stringToBool(value)
		case "dbAdminUser":
			Database.DbAdminUser = value
		case "dbAdminPassword":
			Database.DbAdminPassword = value
		case "dbTimeout":
			Database.DbTimeout = value
		case "mainDb":
			Database.MainDb = value
		}
	}
}

//Load Notification configuration options
func setNotificationConfig(notifSection *configparser.Section) {
	for key, value := range notifSection.Options() {
		switch key {
		case "templateDirectory":
			Notifications.TemplateDir = value
		case "useHtmlTemplates":
			Notifications.UseHtmlTemplates = stringToBool(value)
		case "mainSiteUrl":
			Notifications.MainSiteUrl = value
		case "smtpServer":
			Notifications.SmtpServer = value
		case "useSSL":
			Notifications.UseSSL = stringToBool(value)
		case "smtpPort":
			setIntVal(value, &Notifications.SmtpPort)
		case "smtpUser":
			Notifications.SmtpUser = value
		case "smtpPassword":
			Notifications.SmtpPassword = value
		case "fromEmail":
			Notifications.FromEmail = value
		}
	}
}

func setIntVal(str string, to *int) {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	} else {
		*to = i
	}
}

func setUint64Val(str string, to *uint64) {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		log.Fatal(err)
	} else {
		*to = i
	}
}

func stringToBool(str string) bool {
	if str == "true" {
		return true
	} else {
		return false
	}
}

// Load Logging configuration
func setLogConfig(logSection *configparser.Section) {
	for key, value := range logSection.Options() {
		switch key {
		case "logFile":
			Logger.LogFile = value
		case "maxSize":
			setIntVal(value, &Logger.MaxSize)
		case "maxBackups":
			setIntVal(value, &Logger.MaxBackups)
		case "maxAge":
			setIntVal(value, &Logger.MaxAge)
		}
	}
}

// Load Auth configuration
func setAuthConfig(authSection *configparser.Section) {
	for key, value := range authSection.Options() {
		switch key {
		case "authenticator":
			Auth.Authenticator = strings.ToLower(value)
		case "sessionTimeout":
			setUint64Val(value, &Auth.SessionTimeout)
		case "persistentSessions":
			Auth.PersistentSessions = stringToBool(value)
		case "allowGuestAccess":
			Auth.AllowGuest = stringToBool(value)
		case "allowNewUserRegistration":
			Auth.AllowNewUserRegistration = stringToBool(value)
		case "minPasswordLength":
			setIntVal(value, &Auth.MinPasswordLength)
		}
	}
}

// Load Registry configuration
func setRegistryConfig(registrySection *configparser.Section) {
	for key, value := range registrySection.Options() {
		switch key {
		case "entryTTL":
			setUint64Val(value, &ServiceRegistry.EntryTTL)
		case "cacheRefreshInterval":
			setUint64Val(value, &ServiceRegistry.CacheRefreshInterval)
		}
	}
}

// Load Users configuraiton
func setUsersConfig(userSection *configparser.Section) {
	for key, value := range userSection.Options() {
		switch key {
		case "avatarDB":
			Users.AvatarDb = value
		}
	}
}
