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

package auth_service

import (
	"encoding/json"
	etcd "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/coreos/etcd/client"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/golang.org/x/net/context"
	. "github.com/rhinoman/wikifeat/common/auth"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/common/util"
	"time"
)

// Authentication Manager

type AuthManager struct{}

var sessionsLocation = registry.EtcdPrefix + "/sessions/"

//Get an autheticator
func (am *AuthManager) getAuthenticator(authType string) Authenticator {
	switch authType {
	default:
		return StandardAuthenticator{}
	}
}

//Create a new session (i.e., login)
func (am *AuthManager) Create(username string,
	password string, authType string) (*Session, error) {
	authenticator := am.getAuthenticator(authType)
	sess, err := authenticator.CreateSession(username, password)
	if err != nil {
		return nil, err
	}
	if err = am.registerSession(sess); err != nil {
		return nil, err
	} else {
		return sess, nil
	}
}

//Destroy a session (i.e., logout)
func (am *AuthManager) Destroy(session *Session) error {
	authType := session.AuthType
	authenticator := am.getAuthenticator(authType)
	return authenticator.DestroySession(session.Id)
}

//Get a session
func (am *AuthManager) GetSession(sessionId string) (*Session, error) {
	return nil, nil
}

//Update session
//generate a new session and return the token
func (am *AuthManager) UpdateSession(sessionId string) (string, error) {
	return "", nil
}

//Store the session to etcd
func (am *AuthManager) registerSession(sess *Session) error {
	kapi := registry.GetEtcdKeyAPI()
	ttl := time.Duration(config.Auth.SessionTimeout) * time.Second
	sessBytes, err := json.Marshal(sess)
	sessStr := string(sessBytes)
	if err != nil {
		return err
	}
	_, err = kapi.Set(context.Background(), sessionsLocation+sess.Id, sessStr,
		&etcd.SetOptions{TTL: ttl})
	if err != nil {
		return err
	}
	return nil
}

func NewSession(user string, authType string) *Session {
	token := util.GenToken()
	createdAt := time.Now().UTC()
	sessionTimeout := time.Duration(config.Auth.SessionTimeout) * time.Second
	expiresAt := createdAt.Add(sessionTimeout)
	return &Session{
		Id:        token,
		AuthType:  authType,
		User:      user,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}
