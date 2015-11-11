#!/usr/bin/env python3

"""
    Update Wikifeat couchdb databases from 0.1.1a to 0.2a
    Note: Requires python3

    Changes:
    1.  Added wikit_comments design document to wiki dbs
"""

import json
import common

wiki_ddoc = 'wikit_comments'

getCommentsForPage = dict()

getCommentsForPage['map'] = """
    function(doc){
        if(doc.type==="comment"){
            emit([doc.owning_page, doc.created_time], doc);
        }
    }
"""

getCommentsForPage['reduce'] = "_count"

args = common.parse_args()
conn = common.get_connection(args.use_ssl, args.couch_server, args.couch_port)

credentials = common.get_credentials(args.adminuser, args.adminpass)
get_headers = common.get_headers(credentials)
put_headers = common.put_headers(credentials)

conn.connect()

conn.request("GET", '/_all_dbs', headers=get_headers)
db_list = common.decode_response(conn.getresponse())
wiki_list = [db for db in db_list if db[0:5] == "wiki_"]


# Update the wiiki dbs
for wiki in wiki_list:
    print("Examining " + wiki)
    # Fetch design doc
    ddoc_uri = '/' + wiki + '/_design/' + wiki_ddoc
    conn.request("GET", ddoc_uri, headers=get_headers)
    resp = conn.getresponse()
    common.decode_response(resp)
    if resp.getcode() == 404:
        print("Updating " + wiki)
        ddoc = dict()
        ddoc['views'] = dict()
        ddoc['views']['getCommentsForPage'] = getCommentsForPage
        req_body = json.dumps(ddoc)
        conn.request("PUT", ddoc_uri, body=req_body, headers=put_headers)
        resp = conn.getresponse()
        common.decode_response(resp)
        if resp.getcode() == 201:
            print("Update successful.")
        else:
            print("Update failed.")
    else:
        print("Design doc already exists.  Doing nothing.")



