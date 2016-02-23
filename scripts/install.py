#!/usr/bin/env python3

"""
    Performs a guided install of Wikifeat on the host system
    Note: Requires python3
"""
import sys, os
import util
import setup, config

welcome_text = """
Wikifeat Installation
=====================
This script performs a guided installation and configuration of Wikifeat
on your system. It is intended for a single node system running only
one instance of each service. For multi-node systems and/or more advanced
configurations, see the Wikifeat documentation.
"""

print(welcome_text)
print("")
print("First, we'll need the location and port number for your CouchDB server")
couchdb_server = input("Enter CouchDB hostname or IP (localhost): ")
couchdb_port = input("Enter CouchDB host port number (5984): ")

couchdb_admin = input("Enter the CouchDB admin username: ")
if couchdb_admin == "":
    print("You must specify a CouchDB admin!")
    sys.exit(-1)
couchdb_admin_pass = input("Enter the CouchDB admin password: ")
if couchdb_admin_pass == "":
    print("You must specify a CouchDB admin password!")
    sys.exit(-1)

wf_home_default = os.path.realpath(os.path.join(os.curdir, os.pardir))
wikifeat_home = input("Enter the Wikifeat installation directory(" + wf_home_default + "): ")
if wikifeat_home == "":
    wikifeat_home = wf_home_default

if couchdb_server == "":
    couchdb_server = "localhost"
if couchdb_port == "":
    couchdb_port = "5984"

domain_name = input("Now, enter the domain/host name for your Wikifeat installation (localhost): ")
if domain_name == "":
    domain_name = "localhost"


main_db = "wikifeat_main_db"
avatar_db = "user_avatars"
couch_params = util.CouchParameters()
couch_params.host = couchdb_server
couch_params.port = couchdb_port
couch_params.adminuser = couchdb_admin
couch_params.adminpass = couchdb_admin_pass

master_params = util.MasterUserParameters()

print("Running database setup...")
setup.main(couch_params, main_db, avatar_db, master_params, wikifeat_home)
print("")
print("Configuring Wikifeat...")
config.main(domain_name, couch_params, wikifeat_home)

