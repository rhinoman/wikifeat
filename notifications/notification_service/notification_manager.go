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

package notification_service

// Manager for Notifications

import (
	"bytes"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/config"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/gopkg.in/gomail.v2"
	htemplate "html/template"
	"log"
	"path/filepath"
	"text/template"
)

type NotificationManager struct{}

func (nm *NotificationManager) Send(template string,
	nReq *NotificationRequest) error {
	//TODO: Validate Email Addresses?
	//Create the templates
	var htmlTemplate *htemplate.Template = nil
	plainTemplate, err := nm.LoadPlaintextTemplate(template)
	if err != nil {
		return err
	}
	if config.Notifications.UseHtmlTemplates {
		htmlTemplate, err = nm.LoadHtmlTemplate(template)
		if err != nil {
			log.Printf("Error loading HTML Template: %v", err)
		}
	}
	//Create the email message
	m := gomail.NewMessage()
	m.SetHeader("From", config.Notifications.FromEmail)
	m.SetHeader("To", nReq.To)
	m.SetHeader("Subject", nReq.Subject)
	//Render the templates
	var plainText bytes.Buffer
	var html bytes.Buffer
	//Set the mainSiteUrl in the data object before feeding it to the template
	nReq.Data["_mainSiteUrl"] = config.Notifications.MainSiteUrl
	err = plainTemplate.Execute(&plainText, *nReq)
	if err != nil {
		return err
	}
	m.SetBody("text/plain", plainText.String())
	if htmlTemplate != nil {
		if err = htmlTemplate.Execute(&html, *nReq); err == nil {
			m.AddAlternative("text/html", html.String())
		} else {
			log.Printf("Error executing HTML Template: %v", err)
		}
	}
	d := gomail.NewPlainDialer(config.Notifications.SmtpServer,
		config.Notifications.SmtpPort,
		config.Notifications.SmtpUser,
		config.Notifications.SmtpPassword)
	d.SSL = config.Notifications.UseSSL
	return d.DialAndSend(m)
}

func (nm *NotificationManager) LoadPlaintextTemplate(tName string) (*template.Template, error) {
	tmplDir := config.Notifications.TemplateDir
	filename := tmplDir + "/" + tName + ".txt"
	return template.ParseFiles(filepath.FromSlash(filename))
}

func (nm *NotificationManager) LoadHtmlTemplate(tName string) (*htemplate.Template, error) {
	tmplDir := config.Notifications.TemplateDir
	filename := tmplDir + "/" + tName + ".html"
	return htemplate.ParseFiles(filepath.FromSlash(filename))
}
