#!/usr/bin/env python3

"""
    Update Wikifeat couchdb databases from 0.1.1a to 0.2a
    Note: Requires python3

    Changes:
    1.  Added wikit_comments design document to wiki dbs
    2.  An adjustment to the main database design document for wiki listings
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

main_db = args.main_db
credentials = common.get_credentials(args.adminuser, args.adminpass)
get_headers = common.get_headers(credentials)
put_headers = common.put_headers(credentials)

conn.connect()

# Update the Main db design document
main_ddoc_uri = '/' + main_db + '/_design/wiki_query'
conn.request("GET", main_ddoc_uri, headers=get_headers)
main_ddoc = common.decode_response(conn.getresponse())
main_ddoc['lists']['userWikiList'] = "function(head,req){var row;var user=req['userCtx']['name'];var userRoles=req['userCtx']['roles'];var response={total_rows:0,offset:0,rows:[]};while(row=getRow()){var wikiName=\"wiki_\"+row.id; if(userRoles.indexOf(wikiName+\":read\") > -1 || userRoles.indexOf(wikiName+\":admin\") > -1 || userRoles.indexOf(wikiName+\":write\") > -1 || userRoles.indexOf(\"admin\") > -1 || userRoles.indexOf(\"master\") > -1 || row.value.allowGuest){response.rows.push(row);}}response.total_rows=response.rows.length;send(toJSON(response));}"
req_body = json.dumps(main_ddoc)
conn.request("PUT", main_ddoc_uri, body=req_body, headers=put_headers)
resp = conn.getresponse()
common.decode_response(resp)
if resp.getcode() == 201:
    print("Main ddoc update successful.")
else:
    print("Main ddoc update failed.")

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



