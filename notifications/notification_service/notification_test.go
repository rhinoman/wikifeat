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

package notification_service_test

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/config"
	"github.com/rhinoman/wikifeat/notifications/notification_service"
	"os"
	"testing"
)

var nm = new(notification_service.NotificationManager)

type TemplateData struct {
	To      string
	From    string
	Subject string
	Data    map[string]string
}

func setup() {
	config.LoadDefaults()
	config.Notifications.TemplateDir = "test_templates"
}

func TestLoadTemplate(t *testing.T) {
	setup()
	data := TemplateData{
		To:      "them@otherplace.com",
		From:    "us@ourplace.com",
		Subject: "The Thing",
		Data: map[string]string{
			"user":   "Dude",
			"amount": "Fifty Gazillion",
		},
	}
	//Try the plaintext template
	tmpl, err := nm.LoadPlaintextTemplate("test1")
	if err != nil {
		t.Error(err)
	}
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		t.Error(err)
	}
	//Now try the Html Template
	htmpl, err := nm.LoadHtmlTemplate("test1")
	if err != nil {
		t.Error(err)
	}
	err = htmpl.Execute(os.Stdout, data)
	if err != nil {
		t.Error(err)
	}
}
