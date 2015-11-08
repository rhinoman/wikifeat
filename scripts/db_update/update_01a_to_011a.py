#!/usr/bin/env python3

"""
    Update Wikifeat couchdb databases from 0.1a to 0.1.1a
    Note: Requires python3

    Changes:
    1.  Added reduce functions to various queries in the wiki db design docs
"""

import json
import common

wiki_ddoc = 'wikit'

args = common.parse_args()

# Now do the stuff
conn = common.get_connection(args.use_ssl, args.couch_server, args.couch_port)

credentials = common.get_credentials(args.adminuser, args.adminpass)
get_headers = common.get_headers(credentials)
put_headers = common.put_headers(credentials)

conn.connect()

conn.request("GET", '/_all_dbs', headers=get_headers)
db_list = common.decode_response(conn.getresponse())
wiki_list = [db for db in db_list if db[0:5] == "wiki_"]

# Update the wiki design documents
for wiki in wiki_list:
    print("Updating " + wiki)
    # Fetch design doc
    ddoc_uri = '/' + wiki + '/_design/' + wiki_ddoc
    conn.request("GET", ddoc_uri, headers=get_headers)
    ddoc = common.decode_response(conn.getresponse())
    ddoc['views']['getChildPageIndex']['reduce'] = '_count'
    ddoc['views']['getFileIndex']['reduce'] = '_count'
    ddoc['views']['getHistory']['reduce'] = '_count'
    ddoc['views']['getIndex']['reduce'] = '_count'
    req_body = json.dumps(ddoc)
    conn.request("PUT", ddoc_uri, body=req_body, headers=put_headers)
    resp = conn.getresponse()
    resp_body = common.decode_response(resp)
    if resp.getcode() == 201:
        print("Update successful.")
    else:
        print("Update failed.")

# Last, close the connection
conn.close()