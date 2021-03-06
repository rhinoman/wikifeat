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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/database"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/common/util"
	"net/http"
	"strings"
	"time"
)

type Authenticator interface {
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
	Roles     []string  `json:"roles"`
	AuthType  string    `json:"authType"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserLoginCredentials struct {
	Username string `json:"name"`
	Password string `json:"password"`
	AuthType string `json:"auth_type"`
}

// Implements couchdb-go Auth interface
type WikifeatAuth struct {
	Username  string   `json:"name"`
	Roles     []string `json:"roles"`
	NextToken string   `json:"next_token"`
}

func (wa *WikifeatAuth) AddAuthHeaders(req *http.Request) {
	req.Header.Set("X-Auth-CouchDB-Username", wa.Username)
	rolesString := strings.Join(wa.Roles, ",")
	req.Header.Set("X-Auth-CouchDB-Roles", rolesString)
	// Compute the Auth Token
	secret := database.CouchSecret
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(wa.Username))
	authToken := string(hex.EncodeToString(mac.Sum(nil)))
	req.Header.Set("X-Auth-CouchDB-Token", authToken)
}

func (wa *WikifeatAuth) UpdateAuth(resp *http.Response) {}

func (wa *WikifeatAuth) GetUpdatedAuth() map[string]string {
	ua := make(map[string]string)
	if wa.NextToken != "" {
		ua["AuthSession"] = wa.NextToken
	}
	return ua
}

func (wa *WikifeatAuth) DebugString() string {
	return fmt.Sprintf("Username: %v, Roles: %v", wa.Username, wa.Roles)
}

func AddCsrfCookie(rw http.ResponseWriter, sessToken string) {
	csrfCookie := http.Cookie{
		Name:     "CsrfToken",
		Value:    util.GenHashString(sessToken),
		Path:     "/",
		HttpOnly: false,
	}
	rw.Header().Add("Set-Cookie", csrfCookie.String())
}

func CheckCsrf(request *http.Request) error {
	//Don't bother for Read methods (GET, et. al.,)
	if request.Method == "GET" || request.Method == "HEAD" || request.Method == "OPTIONS" {
		return nil
	}
	//Check CSRF token
	csrfCookie, err := request.Cookie("CsrfToken")
	ourCsrf := request.Header.Get("X-Csrf-Token")
	if err != nil || ourCsrf == "" {
		return UnauthenticatedError()
	} else if token := csrfCookie.Value; token != ourCsrf {
		return UnauthenticatedError()
	} else {
		return nil
	}
}

func GetBasicAuth(req *http.Request) couchdb.Auth {
	//Check first for Basic Auth
	if req.Header.Get("Authorization") != "" {
		return &couchdb.PassThroughAuth{
			AuthHeader: req.Header.Get("Authorization"),
		}
	} else {
		return nil
	}
}

func GetAuth(req *http.Request) (couchdb.Auth, error) {
	//Check if we have a session cookie first
	sessCookie, err := req.Cookie("AuthSession")
	if err != nil {
		return nil, err
	} else if sessCookie.Value == "" {
		return nil, http.ErrNoCookie
	}
	//Attempt to load auth from auth service
	authEndpoint, err := registry.GetServiceLocation("auth")
	if err != nil {
		return nil, err
	}
	reqUrl := authEndpoint + "/api/v1/auth"
	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	request.AddCookie(sessCookie)
	request.Header.Add("Accept", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//The response body should contain our Auth object
	auth := WikifeatAuth{}
	if err = util.DecodeJsonData(resp.Body, &auth); err != nil {
		return nil, err
	} else {
		return &auth, nil
	}
}

func SetAuth(rw http.ResponseWriter, ca couchdb.Auth) {
	authData := ca.GetUpdatedAuth()
	if authData == nil {
		return
	}
	if val, ok := authData["AuthSession"]; ok {
		authCookie := http.Cookie{
			Name:     "AuthSession",
			Value:    val,
			Path:     "/",
			HttpOnly: true,
		}
		rw.Header().Add("Set-Cookie", authCookie.String())
		AddCsrfCookie(rw, util.GenHashString(val))
	}
}

func ClearAuth(rw http.ResponseWriter) {
	//Clear the session cookie
	theCookie := http.Cookie{
		Name:     "AuthSession",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	rw.Header().Add("Set-Cookie", theCookie.String())
}
