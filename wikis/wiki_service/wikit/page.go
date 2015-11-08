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
package wikit

import (
	"time"
)

type History []HistoryViewResult
type PageIndex []PageIndexResult
type FileIndex []FileIndexResult

type PageContent struct {
	Raw       string `json:"raw"`
	Formatted string `json:"formatted"`
}

type Page struct {
	Id          string      `json:"id,omitempty"`
	Slug        string      `json:"slug"`
	DocType     string      `json:"type"`
	Title       string      `json:"title"`
	Owner       string      `json:"owner"`  //a user name
	LastEditor  string      `json:"editor"` //a user name
	Timestamp   time.Time   `json:"timestamp"`
	Content     PageContent `json:"content"`
	Parent      string      `json:"parent"`                    //For page hierarchy: a document id
	Lineage     []string    `json:"lineage"`                   //Parental hierarchy of this page
	OwningPage  string      `json:"owning_page"`               //For page history: a document id
	Attachments []string    `json:"fileAttachments,omitempty"` //A list of file ids
}

type File struct {
	Id          string                `json:"id,omitempty"`
	Name        string                `json:"name"`
	Timestamp   time.Time             `json:"timestamp"`
	UploadedBy  string                `json:"uploadedBy"`
	Description string                `json:"description"`
	DocType     string                `json:"type"`
	Attachments map[string]Attachment `json:"_attachments,omitempty"`
}

type Attachment struct {
	MimeType string `json:"content_type"`
	Digest   string `json:"digest,omitempty"`
	Length   int    `json:"length"`
	RevPos   int    `json:"revpos,omitempty"`
	Stub     bool   `json:"stub,omitempty"`
}

type HistoryEntry struct {
	Editor      string `json:"editor"`
	ContentSize int    `json:"contentSize"`
	DocumentId  string `json:"documentId"`
	DocumentRev string `json:"documentRev"`
}

type ViewResponse struct {
	TotalRows int `json:"total_rows"`
	Offset    int `json:"offset"`
}

type MultiPageResponse struct {
	ViewResponse
	Rows []MultiPageRow `json:"rows,omitempty"`
}

type MultiPageRow struct {
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value struct {
		Rev string `json:"rev"`
	}
	Doc Page `json:"doc"`
}

type HistoryViewResponse struct {
	ViewResponse
	Rows []HistoryViewResult `json:"rows,omitempty"`
}

type HistoryViewResult struct {
	Id    string   `json:"id"`
	Key   []string `json:"key"`
	Value HistoryEntry
}

type SlugViewResponse struct {
	ViewResponse
	Rows []SlugViewResult `json:"rows,omitempty"`
}

type SlugViewResult struct {
	Id    string        `json:"id"`
	Key   string        `json:"key"`
	Value SlugViewEntry `json:"value"`
}

type SlugViewEntry struct {
	Rev  string `json:"pageRev"`
	Page Page   `json:"page"`
}

type KVResponse struct {
	Rows []KvItem `json:"rows"`
}

type KvItem struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type FileIndexViewResponse struct {
	ViewResponse
	Rows []FileIndexResult `json:"rows,omitempty"`
}

type FileIndexResult struct {
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value File   `json:"value"`
}

type PageIndexViewResponse struct {
	ViewResponse
	Rows []PageIndexResult `json:"rows,omitempty"`
}

type PageIndexResult struct {
	Id    string         `json:"id"`
	Key   string         `json:"key"`
	Value PageIndexEntry `json:"value"`
}

type PageIndexEntry struct {
	Id        string    `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Owner     string    `json:"owner"`
	Editor    string    `json:"editor"`
	Timestamp time.Time `json:"timestamp"`
}

type PageViewResult struct {
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value PageValue
}

type PageValue struct {
	Rev string `json:"rev"`
	Pg  Page   `json:"page"`
}

// Page comments
type Comment struct {
	Id            string      `json:"id"`
	DocType       string      `json:"type"`
	OwningPage    string      `json:"owning_page"`
	ParentComment string      `json:"parent_comment"`
	Author        string      `json:"author"`
	Deleted       bool        `json:"deleted"`
	CreatedTime   time.Time   `json:"created_time"`
	ModifiedTime  time.Time   `json:"modified_time"`
	Content       PageContent `json:"content"`
}

type CommentIndexViewResponse struct {
	ViewResponse
	Rows []CommentIndexResult `json:"rows,omitempty"`
}

type CommentIndexResult struct {
	Id    string   `json:"id"`
	Key   []string `json:"key"`
	Value Comment  `json:"value"`
}
