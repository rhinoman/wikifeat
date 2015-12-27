#!/usr/bin/env python3

"""
    Update Wikifeat couchdb databases from 0.3a to 0.3.1a
    Note: Requires python3

    Changes:
    1.  Added by name search to user query design document
"""

import json
import common
import sys

user_ddoc = 'user_queries'

usersByName = dict()
usersByName['map'] = """
function(doc){
    if(doc.type===\"user\" && doc.userPublic){
        emit(doc.name, {name: doc.name, roles: doc.roles, userPublic: doc.userPublic});
        if(doc.userPublic.lastName && doc.name !== doc.userPublic.lastName){
            emit(doc.userPublic.lastName, {name: doc.name, roles: doc.roles, userPublic: doc.userPublic});
        }
        if(doc.userPublic.firstName && doc.name !== doc.userPublic.firstName && doc.userPublic.lastName !== doc.userPublic.firstName){
            emit(doc.userPublic.firstName, {name: doc.name, roles: doc.roles, userPublic: doc.userPublic});
        }
    }
}
"""

usersByName['reduce'] = "_count"

args = common.parse_args()
conn = common.get_connection(args.use_ssl, args.couch_server, args.couch_port)

credentials = common.get_credentials(args.adminuser, args.adminpass)
get_headers = common.get_headers(credentials)
put_headers = common.put_headers(credentials)

conn.connect()

#Update the _users design document
ddoc_uri = '/_users/_design/' + user_ddoc
conn.request("GET", ddoc_uri, headers=get_headers)
resp = conn.getresponse()
if resp.getcode() != 200:
    print("Fetch of user design doc failed")
    sys.exit(1)
ddoc = common.decode_response(resp)
ddoc['views']['usersByName'] = usersByName
req_body = json.dumps(ddoc)
conn.request("PUT", ddoc_uri, body=req_body, headers=put_headers)
resp = conn.getresponse()
resp_body = common.decode_response(resp)

if resp.getcode() == 201:
    print("Update of user design doc successful")
else:
    print("Update of the user design doc failed.")
