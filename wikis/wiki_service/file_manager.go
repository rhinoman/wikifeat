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

package wiki_service

import (
	. "github.com/rhinoman/wikifeat/common/entities"
	. "github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/wikis/wiki_service/wikit"
	"io"
)

type FileManager struct{}

//Gets a list of all 'files' in a wiki
func (fm *FileManager) Index(wiki string, pageNum int, numPerPage int,
	curUser *CurrentUserInfo) (*wikit.FileIndexViewResponse, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.GetFileIndex(pageNum, numPerPage)
}

//Saves a File Record (not the attachment)
func (fm *FileManager) SaveFileRecord(wiki string, file *wikit.File,
	id string, rev string, curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	uploadedBy := curUser.User.UserName
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.SaveFileRecord(file, id, rev, uploadedBy)
}

//Reads a File Record (but not the attachment)
func (fm *FileManager) ReadFileRecord(wiki string,
	file *wikit.File, id string, curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.GetFileRecord(id, file)
}

//Deletes a File Record
func (fm *FileManager) DeleteFile(wiki string,
	id string, curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	theFile := wikit.File{}
	if rev, err := theWiki.GetFileRecord(id, &theFile); err != nil {
		return "", err
	} else {
		return theWiki.DeleteFileRecord(id, rev)
	}
}

//Saves a File's Attachment
func (fm *FileManager) SaveFileAttachment(wiki, id, rev, attName, attType string,
	attContent io.Reader, curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.SaveFileAttachment(id, rev, attName, attType, attContent)
}

//Get file attachment
func (fm *FileManager) GetFileAttachment(wiki, id, rev,
	attType, attName string, curUser *CurrentUserInfo) (io.ReadCloser, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.GetFileAttachment(id, rev, attType, attName)
}
