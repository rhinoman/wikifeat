#!/bin/sh
# Setup script
# 
# Sets up the user database in couchdb 
# Format: setup.sh hostname portnumber

if [ $1 -n ]
then
	echo
	echo "Sets up initial couchdb configuration."
	echo "specify the hostname and port of your couhdb server."
	echo "format is setup.sh hostname port."
	echo
	echo "Assuming localhost:5984"
	echo
fi

host=${1:-localhost}
port=${2:-5984}

echo -n "Enter couchdb admin username: "
read admin_user

echo -n "Enter couchdb admin password: " 
stty -echo
read admin_password
stty echo
echo

#Set the user public fields
url=http://$host:$port/_config/couch_httpd_auth/public_fields
content_type="Content-Type: application/json"
field=userPublic
auth=$admin_user:$admin_password
public_user_fields_cmd="curl -X PUT $url -H \"$content_type\" -d \"$field\" -u $auth"
$public_user_fields_cmd

#Set up views in the user database
url=http://$host:$port/_users/_design/user_queries
echo $userDdoc
writeDesignDocCmd="curl -X PUT $url -H \"$content_type\" --data-binary @user_ddoc.json -u $auth"
echo $writeDesignDocCmd
$writeDesignDocCmd

