package couchdb_test

import "testing"

func TestBulkDocumentClosed(t *testing.T) {
	var err error
	dbName := createTestDb(t)
	conn := getConnection(t)
	db := conn.SelectDB(dbName, adminAuth)
	bulk := db.NewBulkDocument()

	theDoc := TestDocument{
		Title: "My Document",
		Note:  "This is my note",
	}
	theID := getUuid()
	err = bulk.Save(theDoc, theID, "")
	errorify(t, err)

	// first time
	_, err = bulk.Commit()
	errorify(t, err)

	// second times
	_, err = bulk.Commit()
	if err == nil {
		t.Log("ERROR: Must be caused exception when Commit() for second times")
	}

	deleteTestDb(t, dbName)
}

func TestBulkDocumentInsertUpdateDelete(t *testing.T) {
	var err error
	dbName := createTestDb(t)
	conn := getConnection(t)
	db := conn.SelectDB(dbName, adminAuth)

	theID1 := getUuid()
	theRev1 := ""
	theDoc1 := TestDocument{
		Title: "My Document " + theID1,
		Note:  "This is my note",
	}
	theID2 := getUuid()
	theRev2 := ""
	theDoc2 := TestDocument{
		Title: "My Document " + theID2,
		Note:  "This is my note",
	}

	bulkInsert := db.NewBulkDocument()
	okInsertTheDoc1 := false
	err = bulkInsert.Save(theDoc1, theID1, "")
	errorify(t, err)
	okInsertTheDoc2 := false
	err = bulkInsert.Save(theDoc2, theID2, "")
	errorify(t, err)
	insertResults, err := bulkInsert.Commit()
	errorify(t, err)
	for _, insertResult := range insertResults {
		if insertResult.ID == theID1 {
			theRev1 = insertResult.Revision
			okInsertTheDoc1 = insertResult.Ok
		} else if insertResult.ID == theID2 {
			theRev2 = insertResult.Revision
			okInsertTheDoc2 = insertResult.Ok
		}
	}
	if !(okInsertTheDoc1 && okInsertTheDoc2) {
		t.Log("ERROR: failed to insert documents")
	}

	bulkUpdate := db.NewBulkDocument()
	okUpdateTheDoc1 := false
	theDoc1.Note = theDoc1.Note + " " + theRev1
	err = bulkUpdate.Save(theDoc1, theID1, theRev1)
	errorify(t, err)
	okUpdateTheDoc2 := false
	theDoc2.Note = theDoc2.Note + " " + theRev2
	err = bulkUpdate.Save(theDoc2, theID2, theRev2)
	errorify(t, err)
	updateResults, err := bulkUpdate.Commit()
	errorify(t, err)
	for _, updateResult := range updateResults {
		if updateResult.ID == theID1 {
			theRev1 = updateResult.Revision
			okUpdateTheDoc1 = updateResult.Ok
		} else if updateResult.ID == theID2 {
			theRev2 = updateResult.Revision
			okUpdateTheDoc2 = updateResult.Ok
		}
	}
	if !(okUpdateTheDoc1 && okUpdateTheDoc2) {
		t.Log("ERROR: failed to update documents")
	}

	bulkDelete := db.NewBulkDocument()
	okDeleteTheDoc1 := false
	err = bulkDelete.Delete(theID1, theRev1)
	errorify(t, err)
	okDeleteTheDoc2 := false
	err = bulkDelete.Delete(theID2, theRev2)
	errorify(t, err)
	deleteResults, err := bulkDelete.Commit()
	errorify(t, err)
	for _, deleteResult := range deleteResults {
		if deleteResult.ID == theID1 {
			theRev1 = deleteResult.Revision
			okDeleteTheDoc1 = deleteResult.Ok
		} else if deleteResult.ID == theID2 {
			theRev2 = deleteResult.Revision
			okDeleteTheDoc2 = deleteResult.Ok
		}
	}
	if !(okDeleteTheDoc1 && okDeleteTheDoc2) {
		t.Log("ERROR: failed to delete documents")
	}

	deleteTestDb(t, dbName)
}
