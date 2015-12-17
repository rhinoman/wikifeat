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

package user_service

import (
	"bytes"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/nfnt/resize"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	. "github.com/rhinoman/wikifeat/common/entities"
	. "github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/common/util"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"strings"
	"time"
)

type UserAvatarManager struct{}

//Save User Avatar Record
func (uam *UserAvatarManager) Save(id string, rev string,
	avatar *UserAvatar, curUser *CurrentUserInfo) (string, error) {
	nowTime := time.Now().UTC()
	if rev == "" {
		avatar.CreatedAt = nowTime
	}
	avatar.ModifiedAt = nowTime
	var auth couchdb.Auth
	//check for admin
	if util.HasRole(curUser.Roles, AdminRole(MainDbName())) ||
		util.HasRole(curUser.Roles, MasterRole()) {
		auth = AdminAuth
	} else {
		auth = curUser.Auth
	}
	avatarDb := Connection.SelectDB(AvatarDbName(), auth)
	return avatarDb.Save(avatar, id, rev)
}

//Read User Avatar Record
func (uam *UserAvatarManager) Read(id string, avatar *UserAvatar,
	curUser *CurrentUserInfo) (string, error) {
	avatarDb := Connection.SelectDB(AvatarDbName(), curUser.Auth)
	rev, err := avatarDb.Read(id, avatar, nil)
	if err != nil && strings.Contains(err.Error(), ":404:") {
		//No avatar exists, so create it.
		//We need the Admin user for this
		avatarAdminDb := Connection.SelectDB(AvatarDbName(), AdminAuth)
		nowTime := time.Now().UTC()
		avatar := UserAvatar{
			UserName:   id,
			CreatedAt:  nowTime,
			ModifiedAt: nowTime,
		}
		return avatarAdminDb.Save(&avatar, id, "")
	} else {
		return rev, err
	}
}

//Delete a User Avatar Record
func (uam *UserAvatarManager) Delete(id string, curUser *CurrentUserInfo) (string, error) {
	theUser := curUser.User
	var auth couchdb.Auth
	if util.HasRole(theUser.Roles, AdminRole(MainDbName())) ||
		util.HasRole(theUser.Roles, MasterRole()) {
		auth = AdminAuth
	} else {
		auth = curUser.Auth
	}
	avatarDb := Connection.SelectDB(AvatarDbName(), auth)
	//Fetch the record
	avatarRecord := new(UserAvatar)
	rev, err := avatarDb.Read(id, avatarRecord, nil)
	if err != nil {
		return "", err
	}
	return avatarDb.Delete(id, rev)
}

//Save a User Avatar Image
func (uam *UserAvatarManager) SaveImage(id string, rev string, attType string,
	data io.Reader, curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	// Decode the image
	image, _, err := image.Decode(data)
	if err != nil {
		return "", err
	}
	// We need two image sizes, 200px, and a 32px thumbnail
	largeSize := resize.Resize(200, 0, image, resize.Bicubic)
	lRev, err := uam.saveImage(id, rev, "largeSize", largeSize, auth)
	if err != nil {
		return "", err
	}
	thumbnail := resize.Thumbnail(32, 32, image, resize.Bicubic)
	tRev, err := uam.saveImage(id, lRev, "thumbnail", thumbnail, auth)
	if err != nil {
		return "", err
	}
	um := new(UserManager)
	user := User{}
	uRev, err := um.Read(id, &user, curUser)
	if err != nil {
		return "", err
	}
	theUri := ApiPrefix() + "/users" + avatarUri
	user.Public.Avatar = strings.Replace(theUri+"/image", "{user-id}", id, 1)
	user.Public.AvatarThumbnail = strings.Replace(theUri+"/thumbnail", "{user-id}", id, 1)
	if uRev, err = um.Update(id, uRev, &user, curUser); err != nil {
		return "", err
	} else {
		return tRev, nil
	}
}

//Saves an image to the database
func (uam *UserAvatarManager) saveImage(id string, rev string, attName string,
	img image.Image, auth couchdb.Auth) (string, error) {
	// Create a buffer to hold the encoded jpeg
	var buf bytes.Buffer
	// Encode as jpeg
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return "", err
	}
	db := Connection.SelectDB(AvatarDbName(), auth)
	return db.SaveAttachment(id, rev, attName, "image/jpeg", &buf)
}

//Get an Avatar (Large) Image
func (uam *UserAvatarManager) GetLargeAvatar(id string) (io.ReadCloser, error) {
	return uam.getImage(id, "largeSize", AdminAuth)
}

//Get an Avatar (Thumbnail) Image
func (uam *UserAvatarManager) GetThumbnailAvatar(id string) (io.ReadCloser, error) {
	return uam.getImage(id, "thumbnail", AdminAuth)
}

//Fetch image data from database
func (uam *UserAvatarManager) getImage(id string, attName string,
	auth couchdb.Auth) (io.ReadCloser, error) {
	db := Connection.SelectDB(AvatarDbName(), auth)
	return db.GetAttachment(id, "", "image/jpeg", attName)
}
