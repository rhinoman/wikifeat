couchdb-go
==========

[![Build Status](https://travis-ci.org/rhinoman/couchdb-go.svg?branch=master)](https://travis-ci.org/rhinoman/couchdb-go)

Description
-----------

This is my golang CouchDB driver.  There are many like it, but this one is mine.


Installation
------------

```
go get github.com/rhinoman/couchdb-go
```

Documentation
-------------

See the Godoc: http://godoc.org/github.com/rhinoman/couchdb-go

Example Usage
-------------

Connect to a server and create a new document:

```go

type TestDocument struct {
	Title string
	Note string
}

...

var timeout = time.Duration(500 * time.Millisecond)
conn, err := couchdb.NewConnection("127.0.0.1",5984,timeout)
auth := couchdb.BasicAuth{Username: "user", Password: "password" }
db := conn.SelectDB("myDatabase", &auth)

theDoc := TestDocument{
	Title: "My Document",
	Note: "This is a note",
}

theId := genUuid() //use whatever method you like to generate a uuid
//The third argument here would be a revision, if you were updating an existing document
rev, err := db.Save(theDoc, theId, "")  
//If all is well, rev should contain the revision of the newly created
//or updated Document
```








