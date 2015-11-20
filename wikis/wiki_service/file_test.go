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
