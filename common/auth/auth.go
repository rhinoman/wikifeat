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
package auth

import (
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"net/http"
	"time"
)

type Authenticator interface {
	// Reads Auth information/headers from a request and returns a CouchDB auth
	GetAuth(*http.Request) (couchdb.Auth, error)
	// Updates authentication information in the http response before sending it down
	// Usually, this is used to update Auth cookies, CSRF tokens, etc.
	SetAuth(http.ResponseWriter, couchdb.Auth)
	// Creates a new session (i.e., logs a user in)
	CreateSession(string, string) (*Session, error)
	// Destroys a session (i.e., logs a user out)
	DestroySession(string) error
}

type AuthError struct {
	ErrorCode int
	Reason    string
}

func UnauthenticatedError() error {
	return &AuthError{
		ErrorCode: 401,
		Reason:    "Invalid username or password",
	}
}

func (err *AuthError) Error() string {
	return fmt.Sprintf("[Error]:%v: %v", err.ErrorCode, err.Reason)
}

type Session struct {
	Id        string    `json:"id"`
	User      string    `json:"user"`
	AuthType  string    `json:"authType"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserLoginCredentials struct {
	Username string `json:"name"`
	Password string `json:"password"`
	AuthType string `json:"auth_type"`
}
