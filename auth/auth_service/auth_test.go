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

package auth_service_test

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/auth/auth_service"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/database"
	"github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/users/user_service"
	"net/http"
	"testing"
	"time"
)

var timeout = time.Duration(500 * time.Millisecond)
var um = new(user_service.UserManager)
var user = entities.User{
	UserName: "John.Smith",
	Password: "password",
}

func setup(t *testing.T) {
	config.LoadDefaults()
	config.ServiceRegistry.CacheRefreshInterval = 1000
	database.InitDb()
	//This will cause the registry manager to complain, but we don't
	//really need the service being registered here.
	registry.Init("TestAuth", "/database/test/auth")
	//We need to create a user in order to have any sessions, so
	registration := user_service.Registration{
		NewUser: user,
	}
	_, err := um.SetUp(&registration)
	if err != nil {
		t.Error(err)
	}
}

func afterTest(t *testing.T) {
	auth := &couchdb.BasicAuth{
		Username: "John.Smith",
		Password: "password",
	}
	userDoc, _ := database.GetUserFromAuth(auth)
	curUser := &entities.CurrentUserInfo{
		Auth: auth,
		User: userDoc,
	}
	um.Delete(user.UserName, curUser)
	database.DeleteDb(database.MainDbName())
}

func TestSessions(t *testing.T) {
	setup(t)
	defer afterTest(t)
	// Test Standard
	am := auth_service.AuthManager{}
	sess, err := am.Create("John.Smith", "password", "standard")
	if err != nil {
		t.Error(err)
	}
	if sess.User != "John.Smith" {
		t.Error("Username in session not set!")
	}
	t.Logf("Session: %v", sess)
	// Now see if we can read the session back out of etcd
	readSession, err := am.ReadSession(sess.Id)
	if err != nil {
		t.Error(err)
	}
	if readSession.Id != sess.Id ||
		readSession.User != sess.User ||
		readSession.AuthType != sess.AuthType ||
		readSession.CreatedAt != sess.CreatedAt ||
		readSession.ExpiresAt != sess.ExpiresAt {
		t.Error("Sessions are not equal!")
	}
	// Test GetAuth
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Error(err)
	}
	authCookie := &http.Cookie{
		Name:     "AuthSession",
		Value:    readSession.Id,
		Path:     "/",
		HttpOnly: true,
	}
	req.AddCookie(authCookie)
	auth, err := am.GetAuth(req)
	if err != nil {
		t.Error(err)
	}
	auth.AddAuthHeaders(req)
	t.Logf("Auth Headers: %v", req.Header)
	//Try to make a request with this auth
	readUser := entities.User{}
	_, err = database.Connection.GetUser("John.Smith", &readUser, auth)
	if err != nil {
		t.Error(err)
	}
	t.Logf("USER: %v", readUser)
	if len(readUser.Roles) == 0 {
		t.Error("Couldn't read user roles")
	}
	// Test UpdateSession
	updatedSession, err := am.UpdateSession(sess.Id)
	if err != nil {
		t.Error(err)
	}
	if updatedSession.User != sess.User {
		t.Error("User of updated session wrong")
	}
	if updatedSession.Id == sess.Id {
		t.Error("Id not updated")
	}
	t.Logf("New Session: %v", updatedSession)

}
