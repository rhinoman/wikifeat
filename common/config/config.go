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

// Stores configuration for the core services
package config

import (
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/alyu/configparser"
	"log"
	"strconv"
)

var Service struct {
	DomainName       string
	NodeId           string
	Port             string
	ApiVersion       string
	RegistryLocation string
	UseSSL           bool
	SSLCertFile      string
	SSLKeyFile       string
}

var Frontend struct {
	WebAppDir string
	PluginDir string
	Homepage  string
}

var Search struct {
	SearchServerLocation string
}

var Database struct {
	DbAddr          string
	DbPort          string
	UseSSL          bool
	DbAdminUser     string
	DbAdminPassword string
	DbTimeout       string
	MainDb          string
}

var Logger struct {
	LogFile    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

var Auth struct {
	SessionTimeout     int
	PersistentSessions bool
	AllowGuest         bool
	MinPasswordLength  int
}

var ServiceRegistry struct {
	EntryTTL             uint64
	CacheRefreshInterval uint64
}

var Users struct {
	AvatarDb string
}

var Notifications struct {
	TemplateDir      string
	UseHtmlTemplates bool
	SmtpServer       string
	UseSSL           bool
	SmtpPort         int
	SmtpUser         string
	SmtpPassword     string
	FromEmail        string
}

// Initialize Default values
func LoadDefaults() {
	Service.DomainName = "127.0.0.1"
	Service.RegistryLocation = "http://127.0.0.1:2379"
	Service.Port = "6000"
	Service.ApiVersion = "v1"
	Service.NodeId = "cs1"
	Service.UseSSL = false
	Frontend.WebAppDir = "web_app/app"
	Frontend.PluginDir = "plugins"
	Frontend.Homepage = ""
	Database.DbAddr = "127.0.0.1"
	Database.DbPort = "5984"
	Database.UseSSL = false
	Database.DbAdminUser = "adminuser"
	Database.DbAdminPassword = "password"
	Database.DbTimeout = "500"
	Database.MainDb = "main_ut"
	Logger.LogFile = "out.log"
	Logger.MaxSize = 10
	Logger.MaxBackups = 3
	Logger.MaxAge = 30
	Auth.SessionTimeout = 600
	Auth.PersistentSessions = true
	Auth.AllowGuest = true
	Auth.MinPasswordLength = 6
	Users.AvatarDb = "avatar_ut"
	Notifications.TemplateDir = "templates"
	Notifications.UseHtmlTemplates = true
	Notifications.SmtpServer = "localhost"
	Notifications.UseSSL = false
	Notifications.SmtpPort = 587
	Notifications.SmtpUser = "user"
	Notifications.SmtpPassword = "password"
	Notifications.FromEmail = "admin@localhost"
}

// Load config values from file
func LoadConfig(filename string) {
	LoadDefaults()
	log.Printf("\nLoading Configuration from %v\n", filename)
	config, err := configparser.Read(filename)
	fmt.Print(config)
	if err != nil {
		log.Fatal(err)
	}
	serviceSection, err := config.Section("Service")
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
	setServiceConfig(serviceSection)
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

// Load Service configuration options
func setServiceConfig(serverSection *configparser.Section) {
	for key, value := range serverSection.Options() {
		switch key {
		case "domainName":
			Service.DomainName = value
		case "port":
			Service.Port = value
		case "nodeId":
			Service.NodeId = value
		case "registryLocation":
			Service.RegistryLocation = value
		case "apiVersion":
			Service.ApiVersion = value
		case "useSSL":
			Service.UseSSL = stringToBool(value)
		case "sslCertFile":
			Service.SSLCertFile = value
		case "sslKeyFile":
			Service.SSLKeyFile = value
		}
	}
}

// Load Frontend configuration options
func SetFrontendConfig(frontendSection *configparser.Section) {
	for key, value := range frontendSection.Options() {
		switch key {
		case "webAppDir":
			Frontend.WebAppDir = value
		case "pluginDir":
			Frontend.PluginDir = value
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
		case "sessionTimeout":
			setIntVal(value, &Auth.SessionTimeout)
		case "persistentSessions":
			Auth.PersistentSessions = stringToBool(value)
		case "allowGuestAccess":
			Auth.AllowGuest = stringToBool(value)
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
