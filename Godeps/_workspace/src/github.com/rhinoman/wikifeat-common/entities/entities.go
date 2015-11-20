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

package entities

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"time"
)

type ContactInfo struct {
	Phone    string `json:"phone,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Email    string `json:"email,omitempty"`
	Location string `json:"location,omitempty"`
}

type UserPublic struct {
	//Public fields
	LastName        string      `json:"lastName"`
	FirstName       string      `json:"firstName"`
	MiddleName      string      `json:"middleName,omitempty"`
	Title           string      `json:"title,omitempty"`
	Contact         ContactInfo `json:"contactInfo"`
	Avatar          string      `json:"avatar"`
	AvatarThumbnail string      `json:"avatarThumbnail"`
}

type User struct {
	Id             string      `json:"id"`
	UserName       string      `json:"name"`
	Password       string      `json:"password,omitempty"`
	Roles          []string    `json:"roles"`
	Type           string      `json:"type"`
	CreatedAt      time.Time   `json:"createdAt,omitempty"`
	ModifiedAt     time.Time   `json:"modifiedAt,omitempty"`
	PassResetToken ActionToken `json:"password_reset,omitempty"`
	//Public fields
	Public UserPublic `json:"userPublic,omitempty"`
}

type ActionToken struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type UserAvatar struct {
	UserName    string                `json:"_id"`
	CreatedAt   time.Time             `json:"createdAt,omitempty"`
	ModifiedAt  time.Time             `json:"modifiedAt,omitempty"`
	Attachments map[string]Attachment `json:"_attachments,omitempty"`
}

type Attachment struct {
	MimeType string `json:"content_type"`
	Digest   string `json:"digest,omitempty"`
	Length   int    `json:"length"`
	RevPos   int    `json:"revpos,omitempty"`
	Stub     bool   `json:"stub,omitempty"`
}

type CurrentUserInfo struct {
	Auth  couchdb.Auth
	Roles []string
	User  *User
}

//Notificaiton request for the Notification Service
type NotificationRequest struct {
	To      string            `json:"to_email"`
	Subject string            `json:"subject"`
	Data    map[string]string `json:"data"`
}

//WikiRecord entries go in the main database
type WikiRecord struct {
	Id          string    `json:"id,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	HomePageId  string    `json:"homePageId,omitempty"`
	AllowGuest  bool      `json:"allowGuest"`
	Type        string    `json:"type"`
}

func (wr WikiRecord) Validate() error {
	if wr.Name == "" || len(wr.Name) > 128 {
		return &couchdb.Error{
			StatusCode: 400,
			Reason:     "Wiki Name is invalid",
		}
	}
	if wr.Description == "" || len(wr.Description) > 256 {
		return &couchdb.Error{
			StatusCode: 400,
			Reason:     "Wiki Description is invalid",
		}
	}
	if wr.Type != "wiki_record" {
		return &couchdb.Error{
			StatusCode: 400,
			Reason:     "Wiki Doc Type not set",
		}
	}
	return nil
}

type CountResponse struct {
	Rows []CountRecord `json:"rows"`
}

type CountRecord struct {
	Key   string `json:"key,omitempty"`
	Value int    `json:"value"`
}
