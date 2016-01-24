#!/usr/bin/env python3

"""
    Update Wikifeat couchdb databases from 0.4.0a to 0.5.0
    Note: Requires python3

    Changes:
    1.  Added getImageFileIndex view to wiki design documents
"""

import json
import common
import sys

wiki_ddoc = 'wikit'

getImageFileIndex = dict()

getImageFileIndex['map'] = """
function(doc){
    if(doc.type==="file"){
        const att=doc._attachments;
        const contentType=att[Object.keys(att)[0]].content_type;
        if(contentType.substring(0,6)==="image/"){
            emit(doc.name,doc);
        }
    }
}
"""

getImageFileIndex['reduce'] = "_count"

args = common.parse_args()
conn = common.get_connection(args.use_ssl, args.couch_server, args.couch_port)

credentials = common.get_credentials(args.adminuser, args.adminpass)
get_headers = common.get_headers(credentials)
put_headers = common.put_headers(credentials)

# Update all the wiki design docs
conn.request("GET", '/_all_dbs', headers=get_headers)
db_list = common.decode_response(conn.getresponse())
wiki_list = [db for db in db_list if db[0:5] == "wiki_"]

# Update the wiki dbs
for wiki in wiki_list:
    print("Examining " + wiki)
    # Fetch design doc
    ddoc_uri = '/' + wiki + '/_design/' + wiki_ddoc
    conn.request("GET", ddoc_uri, headers=get_headers)
    resp = conn.getresponse()
    ddoc = common.decode_response(resp)
    print("Updating " + wiki)
    ddoc['views']['getImageFileIndex'] = getImageFileIndex
    req_body = json.dumps(ddoc)
    conn.request("PUT", ddoc_uri, body=req_body, headers=put_headers)
    resp = conn.getresponse()
    common.decode_response(resp)
    if resp.getcode() == 201 or resp.getcode() == 200:
        print("Update successful.")
    else:
        print("Update failed.")

# Lastly, close the connection
conn.close()
