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

package wiki_service_test

import (
	"bytes"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	"github.com/rhinoman/wikifeat/users/user_service"
	"github.com/rhinoman/wikifeat/wikis/wiki_service/wikit"
	"io/ioutil"
	"testing"
)

func beforeFileTest(t *testing.T) error {
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
		return err
	}
	return nil
}

func TestFileCRUD(t *testing.T) {
	err := beforeFileTest(t)
	jsAuth := &couchdb.BasicAuth{
		Username: "John.Smith",
		Password: "password",
	}
	if err != nil {
		t.Error(err)
	}
	theUser := User{}
	_, err = grabUser("John.Smith", &theUser, jsAuth)
	if err != nil {
		t.Error(err)
	}
	defer afterTest(&theUser)
	//Create a wiki
	curUser := getCurUser(jsAuth)
	wikiId := getUuid()
	wikiRecord := WikiRecord{
		Name:        "Cafe Project",
		Description: "Wiki for the Cafe Project",
	}
	_, err = wm.Create(wikiId, &wikiRecord, curUser)
	if err != nil {
		t.Error(err)
	}
	defer wm.Delete(wikiId, curUser)
	theFile := wikit.File{
		Name:        "TPS Report",
		Description: "Updated TPS Report Cover Sheet",
	}
	fileId := getUuid()
	//Test Save
	rev, err := fm.SaveFileRecord(wikiId, &theFile, fileId, "", curUser)
	if err != nil {
		t.Error(err)
	}
	t.Logf("File Rev: %v", rev)
	fileData := []byte("TPS COVER SHEET REV.5000")
	fileDataReader := bytes.NewReader(fileData)
	//Test Save Attachment
	aRev, err := fm.SaveFileAttachment(wikiId, fileId, rev, "TPS Report", "text/plain",
		fileDataReader, curUser)
	if err != nil {
		t.Error(err)
	}
	t.Logf("File Att Rev: %v", aRev)
	//Test GetIndex
	fileIndex, err := fm.Index(wikiId, 0, 0, curUser)
	if err != nil {
		t.Error(err)
	}
	if len(fileIndex.Rows) != 1 {
		t.Errorf("File Index should be 1!")
	}
	//Test GetAttachment
	readData, err := fm.GetFileAttachment(wikiId, fileId, aRev,
		"text/plain", "TPS Report", curUser)
	if err != nil {
		t.Error(err)
	}
	theBytes, _ := ioutil.ReadAll(readData)
	if string(theBytes[:]) != "TPS COVER SHEET REV.5000" {
		t.Errorf("CONTENT IS WRONG!")
	}
	//Test GetFileRecord
	readFile := wikit.File{}
	rRev, err := fm.ReadFileRecord(wikiId, &readFile, fileId, curUser)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Read Rev: %v", rRev)
	if readFile.Name != "TPS Report" {
		t.Errorf("File Name was wrong!")
	}
	//Test File Delete
	dRev, err := fm.DeleteFile(wikiId, fileId, curUser)
	if err != nil {
		t.Error(err)
	}
	t.Logf("dRev: %v", dRev)
}
