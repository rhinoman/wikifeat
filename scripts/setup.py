#!/usr/bin/env python3

"""
    Performs initial setup on a Wikifeat installation
    Note: Requires python3
"""

import json
import sys
import util

# Set up some values
welcome_text = """
    Wikifeat Setup
    --------------
    This script performs initial setup of the Wikifeat system.
    It will create a few CouchDB databases, set a few configuration options,
    and populate some design documents.
"""

args = util.parse_args()
credentials = util.get_credentials(args.adminuser, args.adminpass)
gh = util.get_headers(credentials)
ph = util.put_headers(credentials)
main_db = args.main_db

write_role = main_db+":write"
admin_role = "admin"

# This here is the validation function to control writing to the main database
validation_func = """
    function(newDoc, oldDoc, userCtx){
        if((userCtx.roles.indexOf("%s") === -1) &&
            (userCtx.roles.indexOf("%s") === -1) &&
            (userCtx.roles.indexOf("master") === -1) &&
            (userCtx.roles.indexOf("_admin") === -1)){
                throw({forbidden: "Not authorized"});
            }
    }
""" % (admin_role, write_role)

auth_doc = dict()
auth_doc["_id"] = "_design/_auth"
auth_doc["validate_doc_update"] = validation_func

# Establish a connection to couchdb
conn = util.get_connection(args.use_ssl, args.couch_server, args.couch_port)
conn.connect()


def setup_main_db():
    # Create the main database
    print("Creating the main Wikifeat database")
    main_db_url = '/' + main_db
    conn.request("PUT", main_db_url, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 201:
        print("Main database created.")
    elif resp.getcode() == 409 or resp.getcode() == 412:
        print("Main database already exists.")
    else:
        print("Error occurred.")
        sys.exit(-1)
    # Save the auth document
    auth_url = main_db_url + '/_design/_auth'
    conn.request("GET", auth_url, headers=gh)
    resp = conn.getresponse()
    addoc = util.decode_response(resp)
    req_body = ""
    if resp.getcode() == 404:
        req_body = json.dumps(auth_doc)
    elif resp.getcode() == 200:
        addoc['validate_doc_update'] = validation_func
        req_body = json.dumps(addoc)
    if len(req_body) > 1:
        conn.request("PUT", auth_url, body=req_body, headers=ph)
        resp = conn.getresponse()
        util.decode_response(resp)
        if resp.getcode() == 201:
            print("Main auth doc successfully updated.")
        else:
            print("Main auth doc update failed.")
    # Now load the main db security document
    sec_url = main_db_url + '/_security'
    main_sec = util.load_json_file("ddoc/main_access.json")
    req_body = json.dumps(main_sec)
    conn.request("PUT", sec_url, body=req_body, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 200 or resp.getcode() == 201:
        print("Main security doc saved successfully.")
    else:
        print("Main security doc save failed.")


def setup_user_db():
    # Now we'll set up the CouchDB user database
    print("\nSetting user database public fields in CouchDB")
    # First, set the user public fields in the CouchDB config
    url = '/_config/couch_httpd_auth/public_fields'
    field = "\"userPublic\""
    conn.request("PUT", url, body=field, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 200 or resp.getcode() == 201:
        print("User config updated successfully")
    else:
        print("User config update failed!")
    # Now, set up some views in the user database
    url = '/_users/_design/_user_queries'
    # Get the user design doc, if it exists
    conn.request("GET", url, headers=gh)
    resp = conn.getresponse()
    old_ddoc = util.decode_response(resp)
    user_ddoc = util.load_json_file("ddoc/user_ddoc.json")
    if resp.getcode() != 404:
        user_ddoc['_rev'] = old_ddoc['_rev']
    req_body = json.dumps(user_ddoc)
    conn.request("PUT", url, body=req_body, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 200 or resp.getcode() == 201:
        print("User design doc saved successfully.")
    else:
        print("User desgin doc save failed.")


def create_master_user():
    # Now, create a master user
    print("\nCreating a Master User")
    if args.m_uname is None:
        args.m_uname = input("Enter Master username: ")
    if args.m_password is None:
        args.m_password = input("Enter Master user password: ")
    if args.m_ln is None:
        args.m_ln = input("Enter Master User's Last Name: ")
    if args.m_fn is None:
        args.m_fn = input("Enter Master User's First Name: ")
    user_data = dict()
    user_data['_id'] = "org.couchdb.user:" + args.m_uname
    user_data['name'] = args.m_uname
    user_data['type'] = "user"
    user_data['roles'] = list({"master", "admin", "all_users"})
    user_data['password'] = args.m_password
    user_data['userPublic'] = dict({"lastName": args.m_ln, "firstName": args.m_fn})

    url = '/_users/org.couchdb.user:%s' % args.m_uname
    req_body = json.dumps(user_data)
    conn.request("PUT", url, body=req_body, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 200 or resp.getcode() == 201:
        print("Master User created successfully.")
    elif resp.getcode() == 409 or resp.getcode() == 412:
        print("Master User already exists!")
    else:
        print("Creating master user failed.")

# DO THE THINGS
print(welcome_text)
setup_main_db()
setup_user_db()
if args.skip_master is False:
    create_master_user()

