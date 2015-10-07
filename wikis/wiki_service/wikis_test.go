/**   Copyright (c) 2014-present James Adam.  All rights reserved.
*
* This file is part of WikiFeat.
*
*     WikiFeat is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 2 of the License, or
* (at your option) any later version.
*
*     WikiFeat is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
*     You should have received a copy of the GNU General Public License
* along with WikiFeat.  If not, see <http://www.gnu.org/licenses/>.
 */

package wiki_service_test

import (
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/config"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/twinj/uuid"
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

func setup() {
	config.LoadDefaults()
	InitDb()
}

func getUuid() string {
	theUuid := uuid.NewV4()
	return uuid.Formatter(theUuid, uuid.Clean)
}

func grabUser(id string, user *User, auth couchdb.Auth) (string, error) {
	curUser := getCurUser(auth)
	return um.Read(id, user, curUser)
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
	setup()
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
	curUser := &CurrentUserInfo{
		User: user,
		Auth: &couchdb.BasicAuth{
			Username: user.UserName,
			Password: "password",
		},
	}
	DeleteDb(MainDbName())
	um.Delete(user.UserName, curUser)

}

func TestWiki(t *testing.T) {
	beforeTest(t)
	jsAuth := &couchdb.BasicAuth{
		Username: "John.Smith",
		Password: "password",
	}
	curUser := getCurUser(jsAuth)

	theUser := User{}
	_, err := um.Read("John.Smith", &theUser, curUser)
	if err != nil {
		t.Error(err)
	}
	defer afterTest(&theUser)
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
