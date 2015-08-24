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
	Id         string    `json:"id"`
	UserName   string    `json:"name"`
	Password   string    `json:"password,omitempty"`
	Roles      []string  `json:"roles"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	ModifiedAt time.Time `json:"modifiedAt,omitempty"`
	//Public fields
	Public UserPublic `json:"userPublic,omitempty"`
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
