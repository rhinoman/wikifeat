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

package wiki_service_test

import (
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/twinj/uuid"
	"github.com/rhinoman/wikifeat/common/config"
	. "github.com/rhinoman/wikifeat/common/database"
	. "github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/users/user_service"
	"github.com/rhinoman/wikifeat/wikis/wiki_service"
	"testing"
	"time"
)

var timeout = time.Duration(500 * time.Millisecond)
var server = "127.0.0.1"
var adminAuth = &couchdb.BasicAuth{Username: "adminuser", Password: "password"}

var pm = new(wiki_service.PageManager)
var wm = new(wiki_service.WikiManager)
var um = new(user_service.UserManager)
var fm = new(wiki_service.FileManager)
var theUser = User{}
var curUser *CurrentUserInfo
var jsAuth *couchdb.BasicAuth

func setup(mainDb string) {
	config.LoadDefaults()
	config.Database.MainDb = mainDb
	InitDb()
	SetupDb()
}

func getUuid() string {
	theUuid := uuid.NewV4()
	return uuid.Formatter(theUuid, uuid.Clean)
}

func grabUser(id string, user *User, auth couchdb.Auth) (string, error) {
	daUser := getCurUser(auth)
	return um.Read(id, user, daUser)
}

func getCurUser(auth couchdb.Auth) *CurrentUserInfo {
	userDoc, err := GetUserFromAuth(auth)
	if err != nil {
		fmt.Printf("\nERROR: %v\n", err)
	}
	return &CurrentUserInfo{
		Auth: auth,
		User: userDoc,
	}

}

func beforeTest(t *testing.T) {
	time.Sleep(1 * time.Second)
	setup("main_wikis_test")
	user := User{
		UserName: "John.Smith",
		Password: "password",
	}
	registration := user_service.Registration{
		NewUser: user,
	}
	_, err := um.SetUp(&registration)
	if err != nil {
		t.Error(err)
	}
}

func afterTest(user *User) {
	time.Sleep(500 * time.Millisecond)
	daUser := &CurrentUserInfo{
		User: user,
		Auth: &couchdb.BasicAuth{
			Username: user.UserName,
			Password: "password",
		},
	}
	DeleteDb(MainDbName())
	um.Delete(user.UserName, daUser)

}

func TestWikiService(t *testing.T){
	beforeTest(t)
	defer afterTest(&theUser)
	jsAuth := &couchdb.BasicAuth{
		Username: "John.Smith",
		Password: "password",
	}
	curUser = getCurUser(jsAuth)
	_, err := um.Read("John.Smith", &theUser, curUser)
	if err != nil {
		t.Error(err)
	}
	doWikiTest(t)
	doPageTest(t)
	doFileTest(t)
}

func doWikiTest(t *testing.T) {
	wikiId := getUuid()
	otherWikiId := getUuid()
	badWikiId := getUuid()
	//Create a wiki record
	wikiRecord := WikiRecord{
		Name:        "Megasoft All Access",
		Description: "Main Corporate Wiki",
		AllowGuest:  true,
	}
	rev, err := wm.Create(wikiId, &wikiRecord, curUser)
	//Try to create a slug conflict
	badWikiRecord := WikiRecord{
		Name:        "Megasoft All Access",
		Description: "Wiki!",
		AllowGuest:  false,
	}
	_, bErr := wm.Create(badWikiId, &badWikiRecord, curUser)
	if bErr == nil {
		t.Error("Shouldn't have been created!")
	}

	if err != nil {
		t.Error(err)
	}
	otherWikiRecord := WikiRecord{
		Name:        "Megasoft Executives",
		Description: "Executives only",
	}
	rev, err = wm.Create(otherWikiId, &otherWikiRecord, curUser)
	if err != nil {
		t.Error(err)
	}
	//Read record
	rWr := new(WikiRecord)
	rev, err = wm.ReadBySlug(wikiRecord.Slug, rWr, curUser)
	if err != nil || rev == "" {
		t.Error(err)
	}
	t.Logf("WikiRecord: %v", *rWr)
	//Update record
	rWr.Description = "Primary Corporate Wiki"
	rWr.AllowGuest = false
	uRev, err := wm.Update(wikiId, rev, rWr, curUser)
	if err != nil || rev == "" {
		t.Error(err)
	}
	t.Logf("Updated WikiRecord Rev: %v", uRev)
	t.Logf("WikiRecord: %v", *rWr)

	//Try to do it wrong
	oRwr := new(WikiRecord)
	rev, err = wm.Read(otherWikiId, oRwr, curUser)
	if err != nil || rev == "" {
		t.Error(err)
	}
	oRwr.Name = "Megasoft All Access"
	_, err = wm.Update(otherWikiId, rev, oRwr, curUser)
	if err == nil {
		t.Error("Should have been an error!")
	}

	//Test List
	wlr := wiki_service.WikiListResponse{}
	err = wm.GetWikiList(1, 1, false, &wlr, curUser)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Wiki List: %v", wlr)
	if wlr.TotalRows != 2 {
		t.Error("Wrong totalrows!")
	}
	if len(wlr.Rows) != 1 {
		t.Errorf("Wrong length: %v", len(wlr.Rows))
	}
	nextWlr := wiki_service.WikiListResponse{}
	err = wm.GetWikiList(2, 1, false, &nextWlr, curUser)
	if err != nil {
		t.Error(err)
	}
	if len(nextWlr.Rows) != 1 {
		t.Errorf("Wrong length: %v", len(nextWlr.Rows))
	}
	t.Logf("Wiki List: %v", nextWlr)
	//Delete Wiki
	err = wm.Delete(wikiId, curUser)
	if err != nil {
		t.Error(err)
	}
	wm.Delete(otherWikiId, curUser)
}
