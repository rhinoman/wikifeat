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
	"time"
)

type UserAvatarManager struct{}

//Save User Avatar Record
func (uam *UserAvatarManager) Save(id string, rev string,
	avatar *UserAvatar, curUser *CurrentUserInfo) (string, error) {
	theUser := curUser.User
	nowTime := time.Now().UTC()
	if rev == "" {
		avatar.CreatedAt = nowTime
	}
	avatar.ModifiedAt = nowTime
	var auth couchdb.Auth
	//check for admin
	if util.HasRole(theUser.Roles, AdminRole(MainDbName())) ||
		util.HasRole(theUser.Roles, MasterRole()) {
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
	return avatarDb.Read(id, avatar, nil)
}

//Delete a User Avatar Record
func (uam *UserAvatarManager) Delete(id string, curUser *CurrentUserInfo) error {
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
	_, err := avatarDb.Read(id, avatarRecord, nil)
	if err != nil {
		return err
	}
	_, err = avatarDb.Delete(id, "")
	return err
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
	thumbnail := resize.Thumbnail(32, 0, image, resize.Bicubic)
	return uam.saveImage(id, lRev, "thumbnail", thumbnail, auth)
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
func (uam *UserAvatarManager) GetLargeAvatar(id string,
	curUser *CurrentUserInfo) (io.ReadCloser, error) {
	return uam.getImage(id, "largeSize", curUser.Auth)
}

//Get an Avatar (Thumbnail) Image
func (uam *UserAvatarManager) GetThumbnailAvatar(id string,
	curUser *CurrentUserInfo) (io.ReadCloser, error) {
	return uam.getImage(id, "thumbnail", curUser.Auth)
}

//Fetch image data from database
func (uam *UserAvatarManager) getImage(id string, attName string,
	auth couchdb.Auth) (io.ReadCloser, error) {
	db := Connection.SelectDB(AvatarDbName(), auth)
	return db.GetAttachment(id, "", "image/jpeg", attName)
}
