#!/usr/bin/env python3

"""
    Performs initial setup on a Wikifeat installation
    Note: Requires python3
"""

import json
import sys, os
import util

# Set up some values
welcome_text = """
Wikifeat Setup
--------------
This script performs initial setup of the Wikifeat system.
It will create a few CouchDB databases and populate some design documents.
"""

wf_dir = os.curdir


def setup_main_db(conn, main_db):

    write_role = main_db + ":write"
    # This here is the validation function to control writing to the main database
    validation_func = """
    function(newDoc, oldDoc, userCtx){
        if((userCtx.roles.indexOf("%s") === -1) &&
            (userCtx.roles.indexOf("admin") === -1) &&
            (userCtx.roles.indexOf("master") === -1) &&
            (userCtx.roles.indexOf("_admin") === -1)){
                throw({forbidden: "Not authorized"}); }
    }
    """ % write_role

    auth_doc = dict()
    auth_doc["_id"] = "_design/_auth"
    auth_doc["validate_doc_update"] = validation_func

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
    main_sec = util.load_json_file(os.path.join(wf_dir, "scripts/ddoc/main_access.json"))
    req_body = json.dumps(main_sec)
    conn.request("PUT", sec_url, body=req_body, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 200 or resp.getcode() == 201:
        print("Main security doc saved successfully.")
    else:
        print("Main security doc save failed.")
    # Now save the main db design doc
    main_ddoc_url = main_db_url + '/_design/wiki_query'
    conn.request("GET", main_ddoc_url, headers=gh)
    resp = conn.getresponse()
    existing_ddoc = util.decode_response(resp)
    main_ddoc = util.load_json_file(os.path.join(wf_dir, "scripts/ddoc/main_ddoc.json"))
    if resp.getcode() == 200:
        # Set the rev so we can update
        print("Main design doc exists.  Updating.")
        main_ddoc['_rev'] = existing_ddoc['_rev']
    req_body = json.dumps(main_ddoc)
    conn.request("PUT", main_ddoc_url, body=req_body, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 200 or resp.getcode() == 201:
        print("Main design doc saved successfully")
    else:
        print("Main design doc save failed")


def setup_user_db(conn):
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
    url = '/_users/_design/user_queries'
    # Get the user design doc, if it exists
    conn.request("GET", url, headers=gh)
    resp = conn.getresponse()
    old_ddoc = util.decode_response(resp)
    user_ddoc = util.load_json_file(os.path.join(wf_dir, "scripts/ddoc/user_ddoc.json"))
    if resp.getcode() != 404:
        user_ddoc['_rev'] = old_ddoc['_rev']
    req_body = json.dumps(user_ddoc)
    conn.request("PUT", url, body=req_body, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 200 or resp.getcode() == 201:
        print("User design doc saved successfully.")
    else:
        print("User design doc save failed.")


def setup_avatar_db(conn, adb):
    print("Creating the user avatar database")
    adb_url = '/' + adb
    conn.request("PUT", adb_url, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 201:
        print("User Avatar database created.")
    elif resp.getcode() == 409 or resp.getcode() == 412:
        print("Avatar database already exists.")
    else:
        print("Error creating avatar database.")
    # Now save the auth document
    auth_url = adb_url + '/_design/_auth'
    conn.request("GET", auth_url, headers=gh)
    resp = conn.getresponse()
    addoc = util.load_json_file(os.path.join(wf_dir, 'scripts/ddoc/avatar_auth.json'))
    addoc_old = util.decode_response(resp)
    if resp.getcode() == 200:
        print("Avatar auth doc already exists.  Updating.")
        addoc['_rev'] = addoc_old['_rev']
    req_body = json.dumps(addoc)
    conn.request("PUT", auth_url, body=req_body, headers=ph)
    resp = conn.getresponse()
    util.decode_response(resp)
    if resp.getcode() == 200 or resp.getcode() == 201:
        print("Avatar auth doc saved successfully.")
    else:
        print("Avatar auth doc save failed.")


def create_master_user(conn, master):
    # Now, create a master user
    print("\nCreating a Master User")
    if master.user is None:
        print("At least one master user is required per Wikifeat installation.\n"
              "Let's create one now.")
        master.user = input("Enter Master username: ")
    if master.password is None:
        master.password = input("Enter Master user password: ")
    if master.lastname is None:
        master.lastname = input("Enter Master User's Last Name: ")
    if master.firstname is None:
        master.firstname = input("Enter Master User's First Name: ")
    user_data = dict()
    user_data['_id'] = "org.couchdb.user:" + master.user
    user_data['name'] = master.user
    user_data['type'] = "user"
    user_data['roles'] = list({"master", "admin", "all_users"})
    user_data['password'] = master.password
    user_data['userPublic'] = dict({"lastName": master.lastname, "firstName": master.firstname})

    url = '/_users/org.couchdb.user:%s' % master.user
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


def main(couch_params, main_db, avatar_db, master_params, wikifeat_home):
    # Set up credentials
    credentials = util.get_credentials(couch_params.adminuser, couch_params.adminpass)
    global gh, ph
    gh = util.get_headers(credentials)
    ph = util.put_headers(credentials)
    global wf_dir
    wf_dir = wikifeat_home
    # Establish a connection to couchdb
    conn = util.get_connection(
        couch_params.use_ssl,
        couch_params.host,
        couch_params.port)
    conn.connect()
    setup_main_db(conn, main_db)
    setup_user_db(conn)
    setup_avatar_db(conn, avatar_db)
    if master_params.skip_master is False:
        create_master_user(conn, master_params)


if __name__ == "__main__":
    print(welcome_text)
    args = util.parse_args()
    main(util.CouchParameters(args),
         args.main_db,
         args.avatar_db,
         util.MasterUserParameters(args),
         args.wikifeat_home)

