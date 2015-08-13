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

package user_service_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/services"
	. "github.com/rhinoman/wikifeat/users/user_service"
	"testing"
)

var tinyJpeg = "ffd8ffe000104a46494600010101004800480000ffdb0043000302020302020303030304030304050805050404050a070706080c0a0c0c0b0a0b0b0d0e12100d0e110e0b0b1016101113141515150c0f171816141812141514ffc2000b080002000201011100ffc40014000100000000000000000000000000000007ffda00080101000000011effc400161001010100000000000000000000000000050604ffda0008010100010502a1a15303ff00ffc4001a100003010101010000000000000000000001020304050021ffda0008010100063f02e966cdd2d99f3474d2728caeca88a1880a003f07bfffc40017100100030000000000000000000000000001001121ffda0008010100013f21a27af3de5800006013ffda0008010100000010ff00ffc4001510010100000000000000000000000000000100ffda0008010100013f103fe8904041f2800002ffd9"

func setAvatarDbSecurity() error {
	db := services.Connection.SelectDB(services.AvatarDb, services.AdminAuth)
	fmt.Println("Fetching Security document")
	sec, err := db.GetSecurity()
	if err != nil {
		return err
	}
	sec.Admins.Roles = []string{services.AdminRole(services.MainDb),
		services.MasterRole()}
	sec.Members.Roles = []string{"all_users", "guest"}
	fmt.Println("Setting Security for Avatar Database")
	return db.SaveSecurity(*sec)
}

func TestUserAvatars(t *testing.T) {
	setup()
	uam := new(UserAvatarManager)
	//Create the user Avatar Db
	err := services.CreateDb(services.AvatarDb)
	if err != nil {
		t.Error(err)
	}
	err = setAvatarDbSecurity()
	if err != nil {
		t.Error(err)
	}
	defer services.DeleteDb(services.AvatarDb)
	defer services.DeleteDb(services.MainDb)
	smithUser := func() *entities.CurrentUserInfo {
		smithAuth := &couchdb.BasicAuth{Username: "Steven.Smith", Password: "jabberwocky"}
		smith, err := services.GetUserFromAuth(smithAuth)
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
		UserName:    "Steven.Smith",
		UseGravatar: false,
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

}
