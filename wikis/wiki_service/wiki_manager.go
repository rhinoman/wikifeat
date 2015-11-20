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

package wiki_service

// Manager for Wiki Records

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/go-slugification"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/util"
	"github.com/rhinoman/wikifeat/wikis/wiki_service/wikit"
	"log"
	"net/url"
	"strconv"
	"time"
)

type WikiRecordListResult struct {
	Id    string     `json:"id,omitempty"`
	Key   string     `json:"key"`
	Value WikiRecord `json:"value"`
}

type WikiListResponse struct {
	ViewResponse
	Rows []WikiRecordListResult `json:"rows,omitempty"`
}

type WikiSlugViewResponse struct {
	ViewResponse
	Rows []WikiSlugViewResult `json:"rows,omitempty"`
}

type WikiSlugViewResult struct {
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value struct {
		Rev        string     `json:"wikiRev"`
		WikiRecord WikiRecord `json:"wiki_record"`
	} `json:"value"`
}

type CheckSlugResponse struct {
	Rows []KvItem `json:"rows"`
}

type KvItem struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type WikiManager struct{}

func WikiDbName(id string) string {
	return "wiki_" + id
}

//Create a new wiki
func (wm *WikiManager) Create(id string, wr *WikiRecord,
	curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	theUser := curUser.User
	//Verify user is authorized to create wikis
	mainDb := MainDbName()
	if !util.HasRole(theUser.Roles, AdminRole(mainDb)) &&
		!util.HasRole(theUser.Roles, WriteRole(mainDb)) {
		return "", NotAdminError()
	}
	owner := theUser.UserName
	//Add an entry for this wiki to the Main db
	wr.Id = id
	wr.Slug = slugification.Slugify(wr.Name)
	wr.CreatedAt = time.Now().UTC()
	wr.ModifiedAt = time.Now().UTC()
	wr.Type = "wiki_record"
	if err := wr.Validate(); err != nil {
		return "", err
	}
	cDb := Connection.SelectDB(mainDb, auth)
	log.Printf("Adding wiki entry %v to db %v", id, mainDb)
	rev, err := cDb.Save(&wr, id, "")
	//Check for duplicate slugs
	err = wm.checkForDuplicateSlug(wr.Slug)
	if err != nil {
		//Delete the wiki record we just created
		cDb.Delete(id, rev)
		return "", err
	}
	//Create the Wiki
	log.Printf("Creating Wiki: %v", id)
	err = wikit.CreateWiki(Connection, AdminAuth, owner, WikiDbName(id))
	if err != nil {
		//Delete the wiki record from maindb
		cDb.Delete(id, rev)
		return "", err
	}

	//Set Guest Acess
	if err := wm.setGuestAccess(id, wr, auth); err != nil {
		return rev, err
	}
	//wr.Id = id
	return rev, nil
}

// Check for duplicate wiki slug
func (wm *WikiManager) checkForDuplicateSlug(slug string) error {
	params := url.Values{}
	params.Add("key", "\""+slug+"\"")
	params.Add("group", "true")
	response := CheckSlugResponse{}
	mainDb := Connection.SelectDB(MainDbName(), AdminAuth)
	err := mainDb.GetView("wiki_query", "checkUniqueSlug", &response, &params)
	if err != nil {
		return err
	} else if len(response.Rows) <= 0 {
		return nil
	}
	theRecord := response.Rows[0]
	if theRecord.Value > 1 {
		return &couchdb.Error{
			StatusCode: 409,
			Reason:     "Duplicate Wiki slug found",
		}
	} else {
		return nil
	}

}

// Examines the Guest Access flag in the wiki record
// and sets the db security document accordingly
func (wm *WikiManager) setGuestAccess(id string, wr *WikiRecord,
	auth couchdb.Auth) error {
	dbName := WikiDbName(id)
	db := Connection.SelectDB(dbName, auth)
	sec, err := db.GetSecurity()
	if err != nil {
		return err
	}
	members := sec.Members.Users
	if wr.AllowGuest {
		//Make sure guest isn't already a member
		found := false
		for _, member := range members {
			if member == "guest" {
				found = true
				break
			}
		}
		if found == false {
			sec.Members.Users = append(members, "guest")
			if err := db.SaveSecurity(*sec); err != nil {
				return err
			}
		}
		// We need to enable the all_users role.
		// If guests can access, make sure all registered users can too.
		return db.AddRole(AllUsersRole(), false)
	} else {
		found := false
		for i, member := range members {
			if member == "guest" {
				sec.Members.Users = append(members[:i], members[i+1:]...)
				found = true
				break
			}
		}
		if found == true {
			if err := db.SaveSecurity(*sec); err != nil {
				return err
			}
		}
		//remove the all_users role
		return db.RemoveRole(AllUsersRole())
	}

}

//Fetch a wiki record by its slug
func (wm *WikiManager) ReadBySlug(slug string, wikiRecord *WikiRecord,
	curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	mainDb := MainDbName()
	cDb := Connection.SelectDB(mainDb, auth)
	response := WikiSlugViewResponse{}
	err := cDb.GetView("wiki_query", "getWikiBySlug", &response, wikit.SetKey(slug))
	if err != nil {
		return "", err
	}
	if len(response.Rows) > 0 {
		*wikiRecord = response.Rows[0].Value.WikiRecord
		wikiRecord.Id = response.Rows[0].Id
		return response.Rows[0].Value.Rev, nil
	} else {
		return "", NotFoundError()
	}
}

//Fetch a wiki record by its id
func (wm *WikiManager) Read(id string, wikiRecord *WikiRecord,
	curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	mainDb := MainDbName()
	cDb := Connection.SelectDB(mainDb, auth)
	return cDb.Read(id, wikiRecord, nil)
}

//Update a wiki record
func (wm *WikiManager) Update(id string, rev string,
	updateRecord *WikiRecord, curUser *CurrentUserInfo) (string, error) {
	theDb := Connection.SelectDB(MainDbName(), curUser.Auth)
	//Fetch the wiki record
	wr := new(WikiRecord)
	_, err := theDb.Read(id, wr, nil)
	if err != nil {
		return "", err
	}
	//Update select fields
	//Wiki Uuid CANNOT be changed
	//Save the previous data in case we need to undo this (slug conflict, etc.)
	prevDocument := *wr

	//update the data
	wr.Name = updateRecord.Name
	wr.Description = updateRecord.Description
	wr.HomePageId = updateRecord.HomePageId
	wr.AllowGuest = updateRecord.AllowGuest
	wr.ModifiedAt = time.Now().UTC()
	wr.Slug = slugification.Slugify(wr.Name)
	if err = wr.Validate(); err != nil {
		return "", err
	}
	nRev, err := theDb.Save(wr, id, rev)
	if err != nil {
		return "", err
	}
	//Check for duplicate slug
	err = wm.checkForDuplicateSlug(wr.Slug)
	if err != nil {
		//save over with the previous data
		_, sErr := theDb.Save(&prevDocument, id, nRev)
		if sErr != nil {
			return "", sErr
		} else {
			return "", err
		}
	}
	//Update Guest Access
	if err = wm.setGuestAccess(id, wr, curUser.Auth); err != nil {
		return nRev, err
	}
	updateRecord.Id = id
	return nRev, err
}

//Delete a wiki record and the associated wiki database
func (wm *WikiManager) Delete(id string, curUser *CurrentUserInfo) error {
	//Who am I?
	theUser := curUser.User
	auth := curUser.Auth

	mainDb := MainDbName()
	cDb := Connection.SelectDB(mainDb, auth)
	//Check for admin
	if !util.HasRole(theUser.Roles, AdminRole(mainDb)) {
		return NotAdminError()
	}
	//Fetch the wiki record
	wikiRecord := new(WikiRecord)
	_, err := cDb.Read(id, wikiRecord, nil)
	if err != nil {
		return err
	}
	/*if wikiRecord.Id != id {
		return errors.New("WikiRecord doesn't match Database Id")
	}*/
	err = DeleteDb(WikiDbName(id))
	if err != nil {
		return err
	}
	delFunc := func() error {
		rev, err := cDb.Read(id, wikiRecord, nil)
		_, err = cDb.Delete(id, rev)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
	err = util.Retry(3, delFunc)
	if err != nil {
		return err
	}
	return nil
}

//Get list of all wikis
func (wm *WikiManager) GetWikiList(pageNum int, numPerPage int, memberOnly bool,
	wlr *WikiListResponse, curUser *CurrentUserInfo) error {
	params := url.Values{}
	if numPerPage != 0 {
		params.Add("limit", strconv.Itoa(numPerPage))
	}
	skip := numPerPage * (pageNum - 1)
	if skip > 0 {
		params.Add("skip", strconv.Itoa(skip))
	}
	auth := curUser.Auth
	mainDb := MainDbName()
	cDb := Connection.SelectDB(mainDb, auth)
	var err error
	if memberOnly {
		err = cDb.GetList("wiki_query", "userWikiList", "getWikis",
			&wlr, &params)
	} else {
		err = cDb.GetView("wiki_query", "getWikis", &wlr, &params)
	}
	if err != nil {
		return err
	}
	return nil
}
