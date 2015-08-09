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
	//"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	. "github.com/rhinoman/wikifeat/common/entities"
	//. "github.com/rhinoman/wikifeat/common/services"
	"io"
)

type UserAvatarManager struct{}

//Save a User Avatar
func (um *UserManager) Save(id string, rev string,
	data io.Reader, curUser *CurrentUserInfo) (string, error) {
	return "", nil

}
