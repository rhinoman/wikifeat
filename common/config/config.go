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

package config

import (
	"github.com/rhinoman/wikifeat/common/util"
	"log"
	"path"
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
	Authenticator      string
	SessionTimeout     uint64
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
	MainSiteUrl      string
	FromEmail        string
}

// Initialize Default values
func LoadDefaults() {
	execDir, err := util.GetExecDirectory()
	if err != nil {
		log.Fatal(err)
	}
	Service.DomainName = "127.0.0.1"
	Service.RegistryLocation = "http://127.0.0.1:2379"
	Service.Port = "6000"
	Service.ApiVersion = "v1"
	Service.NodeId = "cs1"
	Service.UseSSL = false
	ServiceRegistry.CacheRefreshInterval = 75
	ServiceRegistry.EntryTTL = 60
	Frontend.WebAppDir = path.Join(execDir, "web_app/app")
	Frontend.PluginDir = path.Join(execDir, "plugins")
	Frontend.Homepage = ""
	Database.DbAddr = "127.0.0.1"
	Database.DbPort = "5984"
	Database.UseSSL = false
	Database.DbAdminUser = "adminuser"
	Database.DbAdminPassword = "password"
	Database.DbTimeout = "0"
	Database.MainDb = "main_ut"
	Logger.LogFile = "wikifeat-service.log"
	Logger.MaxSize = 10
	Logger.MaxBackups = 3
	Logger.MaxAge = 30
	Auth.Authenticator = "standard"
	Auth.SessionTimeout = 600
	Auth.PersistentSessions = true
	Auth.AllowGuest = true
	Auth.MinPasswordLength = 6
	Users.AvatarDb = "avatar_ut"
	Notifications.TemplateDir = "templates"
	Notifications.UseHtmlTemplates = true
	Notifications.MainSiteUrl = "http://localhost:8081"
	Notifications.SmtpServer = "localhost"
	Notifications.UseSSL = false
	Notifications.SmtpPort = 587
	Notifications.SmtpUser = "user"
	Notifications.SmtpPassword = "password"
	Notifications.FromEmail = "admin@localhost"
}
