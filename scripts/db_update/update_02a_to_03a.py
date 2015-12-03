#!/usr/bin/env python3

"""
    Update Wikifeat couchdb databases from 0.2a to 0.3a
    Note: Requires python3

    Changes:
    1.  Changes to make API more consistent
"""

import json
import common

wiki_ddoc = 'wikit'
comment_ddoc = 'wikit_comments'

getHistory = dict()
getHistory['map'] = """
    function(doc) {
        if(doc.type==="page"){
            emit([doc.owningPage, doc.timestamp],
            {
                documentId: doc._id,
                documentRev: doc._rev,
                editor: doc.editor,
                contentSize: doc.content.raw.length
            });
        }
    }
"""
getHistory['reduce'] = "_count"

getIndex = dict()
getIndex['map'] = """
    function(doc){
        if(doc.type==="page" && doc._id === doc.owningPage){
            emit(doc.title, {
                id: doc._id,
                slug: doc.slug,
                title: doc.title,
                owner: doc.owner,
                editor: doc.editor,
                timestamp: doc.timestamp
            });
        }
    }
"""
getIndex['reduce'] = "_count"

getChildPageIndex = dict()
getChildPageIndex['map'] = """
    function(doc) {
        if(doc.type==="page" && doc._id === doc.owningPage){
            emit(doc.parent, {
                id: doc._id,
                slug: doc.slug,
                title: doc.title,
                owner: doc.owner,
                editor: doc.editor,
                timestamp: doc.timestamp
            });
        }
    }
"""
getChildPageIndex['reduce'] = "_count"

getCommentsForPage = dict()
getCommentsForPage['map'] = """
    function(doc){
        if(doc.type==="comment"){
            emit([doc.owningPage, doc.createdTime], doc);
        }
    }
"""
getCommentsForPage['reduce'] = "_count"

args = common.parse_args()
conn = common.get_connection(args.use_ssl, args.couch_server, args.couch_port)

main_db = args.main_db
credentials = common.get_credentials(args.adminuser, args.adminpass)
get_headers = common.get_headers(credentials)
put_headers = common.put_headers(credentials)

conn.connect()

#Update wiki design documents
conn.request("GET", '/_all_dbs', headers=get_headers)
db_list = common.decode_response(conn.getresponse())
wiki_list = [db for db in db_list if db[0:5] == "wiki_"]

# Update the wiki dbs
for wiki in wiki_list:
    print("Examining " + wiki)
    # Fetch wiki design doc
    ddoc_uri = '/' + wiki + '/_design/' + wiki_ddoc
    conn.request("GET", ddoc_uri, headers=get_headers)
    resp = conn.getresponse()
    ddoc = common.decode_response(resp)
    ddoc['views']['getHistory'] = getHistory
    ddoc['views']['getIndex'] = getIndex
    ddoc['views']['getChildPageIndex'] = getChildPageIndex
    req_body = json.dumps(ddoc)
    conn.request("PUT", ddoc_uri, body=req_body, headers=put_headers)
    resp = conn.getresponse()
    resp_body = common.decode_response(resp)
    if resp.getcode() == 201:
        print("Update of wiki design doc successful.")
    else:
        print("Update of wiki design doc failed.")
    # Fetch comment design doc
    ddoc_uri = '/' + wiki + '/_design/' + comment_ddoc
    conn.request("GET", ddoc_uri, headers=get_headers)
    resp = conn.getresponse()
    ddoc = common.decode_response(resp)
    ddoc['views']['getCommentsForPage'] = getCommentsForPage
    req_body = json.dumps(ddoc)
    conn.request("PUT", ddoc_uri, body=req_body, headers=put_headers)
    resp = conn.getresponse()
    common.decode_response(resp)
    if resp.getcode() == 201 or resp.getcode() == 200:
        print("Update of Comment Design doc successful.")
    else:
        print("Update of Comment Design doc failed.")
    print("Updating records in " + wiki + ". This might take a while.")
    # Update all the documents with corrected keys
    all_uri = '/' + wiki + '/_all_docs'
    conn.request("GET", all_uri, headers=get_headers)
    resp = conn.getresponse()
    resp_body = common.decode_response(resp)
    doc_count = 0
    for row in resp_body["rows"]:
        doc_id = row['id']
        doc_rev = row['value']['rev']
        # Fetch the individual document
        doc_uri = '/' + wiki + '/' + doc_id
        conn.request("GET", doc_uri, headers=get_headers)
        resp = conn.getresponse()
        doc = common.decode_response(resp)
        modified = False
        if "owning_page" in doc:
                doc["owningPage"] = doc["owning_page"]
                del doc["owning_page"]
                modified = True
        if "type" in doc and doc["type"] == "page":
            if "comments_disabled" in doc:
                doc["commentsDisabled"] = doc["comments_disabled"]
                del doc["comments_disabled"]
                modified = True
        if "type" in doc and doc["type"] == "comment":
            if "created_at" in doc:
                doc["createdAt"] = doc["created_at"]
                del doc["created_at"]
                modified = True
            if "modified_at" in doc:
                doc["modifiedAt"] = doc["modified_at"]
                del doc["modified_at"]
                modified = True
        if modified is True:
            req_body = json.dumps(doc)
            conn.request("PUT", doc_uri, body=req_body, headers=put_headers)
            resp = conn.getresponse()
            resp_body = common.decode_response(resp)
            if resp.getcode() == 201:
                doc_count += 1
    print("Updated " + str(doc_count) + " records in wiki " + wiki)






