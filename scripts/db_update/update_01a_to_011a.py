#!/usr/bin/env python3

"""
    Update Wikifeat couchdb databases from 0.1a to 0.1.1a
    Note: Requires python3

    Changes:
    1.  Added reduce functions to various queries in the wiki db design docs
"""

from argparse import ArgumentParser
import http.client
import json
from base64 import b64encode

wiki_ddoc = 'wikit'

def decode_response(resp):
    return json.loads(resp.read().decode('utf-8'))


# We need the CouchDB admin credentials
# These can be provided as command line arguments or via prompt

parser = ArgumentParser()
# CouchDB URL is required
parser.add_argument('couch_server', type=str)
parser.add_argument('couch_port', type=int)
parser.add_argument('-u', '--user', dest='adminuser',
                    help='CouchDB admin user')
parser.add_argument('-p', '--password', dest='adminpass',
                    help='CouchDB admin password')
# Note: your python must be compiled with SSL support to use HTTPS
parser.add_argument('--use_ssl', dest='use_ssl', action='store_true')
parser.set_defaults(use_ssl=False)

args = parser.parse_args()

if args.adminuser is None:
    args.adminuser = input("Enter CouchDB admin username: ")

if args.adminpass is None:
    args.adminpass = input("Enter CouchDB admin password: ")

credentials = b64encode(bytes(args.adminuser + ':' + args.adminpass, 'utf-8')).decode('utf-8')
get_headers = {
    'Accept': 'application/json',
    'Authorization': 'Basic %s' % credentials
}

put_headers = get_headers.copy()
put_headers.update({'Content-Type': 'application/json'})

# Now do the stuff
conn = None
if args.use_ssl:
    conn = http.client.HTTPSConnection(args.couch_server, args.couch_port)
else:
    conn = http.client.HTTPConnection(args.couch_server, args.couch_port)

conn.connect()

conn.request("GET", '/_all_dbs', headers=get_headers)
db_list = decode_response(conn.getresponse())
wiki_list = [db for db in db_list if db[0:5] == "wiki_"]

# Update the wiki design documents
for wiki in wiki_list:
    print("Updating " + wiki)
    # Fetch design doc
    ddoc_uri = '/' + wiki + '/_design/' + wiki_ddoc
    conn.request("GET", ddoc_uri, headers=get_headers)
    ddoc = decode_response(conn.getresponse())
    ddoc['views']['getChildPageIndex']['reduce'] = '_count'
    ddoc['views']['getFileIndex']['reduce'] = '_count'
    ddoc['views']['getHistory']['reduce'] = '_count'
    ddoc['views']['getIndex']['reduce'] = '_count'
    req_body = json.dumps(ddoc)
    conn.request("PUT", ddoc_uri, body=req_body, headers=put_headers)
    resp = conn.getresponse()
    resp_body = decode_response(resp)
    if resp.getcode() == 201:
        print("Update successful.")
    else:
        print("Update failed.")

# Last, close the connection
conn.close()