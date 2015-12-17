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
