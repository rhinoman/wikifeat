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

package notification_service

// Manager for Notifications

import (
	"bytes"
	"github.com/rhinoman/wikifeat/common/config"
	. "github.com/rhinoman/wikifeat/common/entities"
	"gopkg.in/gomail.v2"
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
