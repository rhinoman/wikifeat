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
	Id              string      `json:"id,omitempty"`
	Slug            string      `json:"slug"`
	DocType         string      `json:"type"`
	Title           string      `json:"title"`
	Owner           string      `json:"owner"`  //a user name
	LastEditor      string      `json:"editor"` //a user name
	Timestamp       time.Time   `json:"timestamp"`
	Content         PageContent `json:"content"`
	Parent          string      `json:"parent"`                    //For page hierarchy: a document id
	Lineage         []string    `json:"lineage"`                   //Parental hierarchy of this page
	OwningPage      string      `json:"owning_page"`               //For page history: a document id
	DisableComments bool        `json:"comments_disabled"`         //disallow comments for this page
	Attachments     []string    `json:"fileAttachments,omitempty"` //A list of file ids
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
	Id           string      `json:"id"`
	Rev          string      `json:"_rev,omitempty"`
	DocType      string      `json:"type"`
	OwningPage   string      `json:"owning_page"`
	Author       string      `json:"author"`
	CreatedTime  time.Time   `json:"created_time"`
	ModifiedTime time.Time   `json:"modified_time"`
	Content      PageContent `json:"content"`
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
