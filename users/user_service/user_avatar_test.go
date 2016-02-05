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
package user_service_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/database"
	"github.com/rhinoman/wikifeat/common/entities"
	. "github.com/rhinoman/wikifeat/users/user_service"
	"io/ioutil"
	"testing"
)

var tinyJpeg = "ffd8ffe000104a46494600010101004800480000ffdb0043000302020302020303030304030304050805050404050a070706080c0a0c0c0b0a0b0b0d0e12100d0e110e0b0b1016101113141515150c0f171816141812141514ffc2000b080002000201011100ffc40014000100000000000000000000000000000007ffda00080101000000011effc400161001010100000000000000000000000000050604ffda0008010100010502a1a15303ff00ffc4001a100003010101010000000000000000000001020304050021ffda0008010100063f02e966cdd2d99f3474d2728caeca88a1880a003f07bfffc40017100100030000000000000000000000000001001121ffda0008010100013f21a27af3de5800006013ffda0008010100000010ff00ffc4001510010100000000000000000000000000000100ffda0008010100013f103fe8904041f2800002ffd9"

func setAvatarDbSecurity() error {
	db := database.Connection.SelectDB(database.AvatarDb, database.AdminAuth)
	fmt.Println("Fetching Security document")
	sec, err := db.GetSecurity()
	if err != nil {
		return err
	}
	sec.Admins.Roles = []string{database.AdminRole(database.MainDb),
		database.MasterRole()}
	sec.Members.Roles = []string{"all_users", "guest"}
	fmt.Println("Setting Security for Avatar Database")
	return db.SaveSecurity(*sec)
}

func TestUserAvatars(t *testing.T) {
	setup()
	uam := new(UserAvatarManager)
	//Create the user Avatar Db
	err := database.CreateDb(database.AvatarDb)
	if err != nil {
		t.Error(err)
	}
	err = setAvatarDbSecurity()
	if err != nil {
		t.Error(err)
	}
	defer database.DeleteDb(database.AvatarDb)
	defer database.DeleteDb(database.MainDb)
	smithUser := func() *entities.CurrentUserInfo {
		smithAuth := &couchdb.BasicAuth{Username: "Steven.Smith", Password: "jabberwocky"}
		smith, err := database.GetUserFromAuth(smithAuth)
		if err != nil {
			t.Error(err)
			return &entities.CurrentUserInfo{}
		}
		return &entities.CurrentUserInfo{
			Auth: smithAuth,
			User: smith,
		}
	}
	//after test cleanup
	defer func() {
		getSmithUser := func() *entities.CurrentUserInfo {
			return smithUser()
		}
		um.Delete("Steven.Smith", getSmithUser())
	}()
	//Register user, get things set up, etc.
	user := entities.User{
		UserName: "Steven.Smith",
		Password: "jabberwocky",
	}
	registration := Registration{
		NewUser: user,
	}
	rev, err := um.SetUp(&registration)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("User Rev: %v", rev)
	}
	//Now, create a User Avatar Record
	uar := entities.UserAvatar{
		UserName: "Steven.Smith",
	}
	aRev, err := uam.Save("Steven.Smith", "", &uar, smithUser())
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("User Avatar Rev: %v", aRev)
	}
	//Save the image
	t.Logf("Decoding Hex String")
	imageBytes, err := hex.DecodeString(tinyJpeg)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Creating io.Reader and sending to UAM")
	imageReader := bytes.NewReader(imageBytes)
	iRev, err := uam.SaveImage("Steven.Smith", aRev, "image/jpeg",
		imageReader, smithUser())
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("User Avatar Image Rev: %v", iRev)
	}
	//Read the Avatar Record
	readRecord := entities.UserAvatar{}
	rRev, err := uam.Read("Steven.Smith", &readRecord, smithUser())
	if err != nil {
		t.Error(err)
	}
	t.Logf("Read Avatar Record with Revision: %v", rRev)
	//Read the Avatar Large Image
	imgData, err := uam.GetLargeAvatar("Steven.Smith")
	if err != nil {
		t.Error(err)
	}
	imgBytes, err := ioutil.ReadAll(imgData)
	imgData.Close()
	if len(imgBytes) == 0 {
		t.Error("Image was zero length!")
	}
	t.Logf("Read Avatar Large Image with %v bytes", len(imgBytes))
	//Read the Avatar Thumbnail
	imgData, err = uam.GetThumbnailAvatar("Steven.Smith")
	if err != nil {
		t.Error(err)
	}
	imgBytes, err = ioutil.ReadAll(imgData)
	imgData.Close()
	if len(imgBytes) == 0 {
		t.Error("Image was zero length!")
	}
	t.Logf("Read Avatar Thumbnail Image with %v bytes", len(imgBytes))
	//Delete Avatar Record
	dRev, err := uam.Delete("Steven.Smith", smithUser())
	if err != nil {
		t.Error(err)
	}
	t.Logf("Delete Avatar with Rev: %v", dRev)
}
