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
	"errors"
	etcd "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/coreos/etcd/client"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/golang.org/x/net/context"
	. "github.com/rhinoman/wikifeat/common/auth"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/common/util"
	"net/http"
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
	if err = am.saveSession(sess); err != nil {
		return nil, err
	} else {
		return sess, nil
	}
}

//Destroy a session (i.e., logout)
func (am *AuthManager) Destroy(session *Session) error {
	authType := session.AuthType
	authenticator := am.getAuthenticator(authType)
	err := authenticator.DestroySession(session.Id)
	if err != nil {
		return err
	}
	//Remove the session node from etcd
	kapi := registry.GetEtcdKeyAPI()
	sessionLocation := sessionsLocation + session.Id
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = kapi.Delete(ctx, sessionLocation, nil)
	return err
}

//Get a session
func (am *AuthManager) ReadSession(sessionId string) (*Session, error) {
	kapi := registry.GetEtcdKeyAPI()
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	sessionLocation := sessionsLocation + sessionId
	resp, err := kapi.Get(ctx, sessionLocation, &etcd.GetOptions{Recursive: false})
	if err != nil {
		return nil, err
	}
	sess, err := am.processResponse(resp)
	if err != nil {
		return nil, err
	} else {
		return sess, nil
	}
}

//Update session
//generate a new sessionId and return the new session
func (am *AuthManager) UpdateSession(sessionId string) (*Session, error) {
	curSession, err := am.ReadSession(sessionId)
	if err != nil {
		return nil, err
	}
	username := curSession.User
	authType := curSession.AuthType
	if username == "" {
		return nil, UnauthenticatedError()
	}
	//Fetch the user object
	user, err := getUser(username)
	if err != nil {
		return nil, UnauthenticatedError()
	}
	userCtx := couchdb.UserContext{
		Name:  username,
		Roles: user.Roles,
	}
	newSession := NewSession(&userCtx, authType)
	err = am.saveSession(newSession)
	if err != nil {
		return nil, err
	} else {
		return newSession, nil
	}
}

// Takes an etcd response and extracts the Session
func (am *AuthManager) processResponse(resp *etcd.Response) (*Session, error) {
	node := resp.Node
	if node.Dir {
		return nil, errors.New("Session node is a directory?!")
	}
	value := []byte(node.Value)
	if len(value) == 0 {
		return nil, errors.New("Session value is empty")
	}
	session := Session{}
	err := json.Unmarshal(value, &session)
	if err != nil {
		return nil, err
	} else {
		return &session, nil
	}
}

//Store the session to etcd
func (am *AuthManager) saveSession(sess *Session) error {
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

//Produce an auth object from the given request
func (am *AuthManager) GetAuth(req *http.Request, authType string) (*couchdb.ProxyAuth, error) {
	authenticator := am.getAuthenticator(authType)
	sessionId, err := authenticator.GetSessionId(req)
	if err != nil {
		return nil, UnauthenticatedError()
	}
	//get the session
	session, err := am.ReadSession(sessionId)
	if err != nil {
		return nil, UnauthenticatedError()
	} else {
		return &couchdb.ProxyAuth{
			Username:  session.User,
			Roles:     session.Roles,
			AuthToken: sessionId,
		}, nil
	}
}

func getUser(username string) (*entities.User, error) {
	user := entities.User{}
	_, err := services.Connection.GetUser(username, &user, services.AdminAuth)
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func NewSession(userCtx *couchdb.UserContext, authType string) *Session {
	token := util.GenToken()
	createdAt := time.Now().UTC()
	sessionTimeout := time.Duration(config.Auth.SessionTimeout) * time.Second
	expiresAt := createdAt.Add(sessionTimeout)
	return &Session{
		Id:        token,
		AuthType:  authType,
		User:      userCtx.Name,
		Roles:     userCtx.Roles,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}
