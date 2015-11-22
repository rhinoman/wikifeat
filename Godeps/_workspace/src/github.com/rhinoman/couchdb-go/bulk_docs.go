package couchdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

type bulkDoc struct {
	_id      string
	_rev     string
	_deleted bool
	doc      interface{}
}

func (b bulkDoc) MarshalJSON() ([]byte, error) {
	out := make(map[string]interface{})
	if !b._deleted {
		bValue := reflect.Indirect(reflect.ValueOf(b.doc))
		bType := bValue.Type()
		for i := 0; i < bType.NumField(); i++ {
			field := bType.Field(i)
			name := bType.Field(i).Name
			jsonKey := field.Tag.Get("json")
			if jsonKey == "" {
				jsonKey = name
			}
			out[jsonKey] = bValue.FieldByName(name).Interface()
		}
	} else {
		out["_deleted"] = true
	}
	out["_id"] = b._id
	if b._rev != "" {
		out["_rev"] = b._rev
	}
	return json.Marshal(out)
}

// BulkDocument Bulk Document API
// http://docs.couchdb.org/en/1.6.1/api/database/bulk-api.html#db-bulk-docs
type BulkDocument struct {
	docs   []bulkDoc
	db     *Database
	closed bool
}

// NewBulkDocument New BulkDocument instance
func (db *Database) NewBulkDocument() *BulkDocument {
	b := &BulkDocument{}
	b.db = db
	return b
}

// Save Save document
func (b *BulkDocument) Save(doc interface{}, id, rev string) error {
	if id == "" {
		return fmt.Errorf("No ID specified")
	}
	b.docs = append(b.docs, bulkDoc{id, rev, false, doc})
	return nil
}

// Delete Delete document
func (b *BulkDocument) Delete(id, rev string) error {
	if id == "" {
		return fmt.Errorf("No ID specified")
	}
	if rev == "" {
		return fmt.Errorf("No Revision specified")
	}
	b.docs = append(b.docs, bulkDoc{id, rev, true, nil})
	return nil
}

// BulkDocumentResult Bulk Document Response
type BulkDocumentResult struct {
	Ok       bool    `json:"ok"`
	ID       string  `json:"id"`
	Revision string  `json:"rev"`
	Error    *string `json:"error"`
	Reason   *string `json:"reason"`
}

func getBulkDocumentResult(resp *http.Response) ([]BulkDocumentResult, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results []BulkDocumentResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Commit POST /{db}/_bulk_docs
func (b *BulkDocument) Commit() ([]BulkDocumentResult, error) {
	if !b.closed {
		b.closed = true
		url, err := buildUrl(b.db.dbName, "_bulk_docs")
		if err != nil {
			return nil, err
		}
		var headers = make(map[string]string)
		headers["Content-Type"] = "application/json"
		headers["Accept"] = "application/json"
		bd := make(map[string]interface{})
		bd["docs"] = b.docs
		data, numBytes, err := encodeData(bd)
		if err != nil {
			return nil, err
		}
		headers["Content-Length"] = strconv.Itoa(numBytes)
		//Yes, this needs to be here.
		//Yes, I know the Golang http.Client doesn't support expect/continue
		//This is here to work around a bug in CouchDB.  It shouldn't work, and yet it does.
		//See: http://stackoverflow.com/questions/30541591/large-put-requests-from-go-to-couchdb
		//Also, I filed a bug report: https://issues.apache.org/jira/browse/COUCHDB-2704
		//Go net/http needs to support the HTTP/1.1 spec, or CouchDB needs to get fixed.
		//If either of those happens in the future, I can revisit this.
		//Unless I forget, which I'm sure I will.
		if numBytes > 4000 {
			headers["Expect"] = "100-continue"
		}
		resp, err := b.db.connection.request("POST", url, data, headers, b.db.auth)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return getBulkDocumentResult(resp)
	}
	return nil, fmt.Errorf("CouchDB: Bulk Document has already been executed")
}
