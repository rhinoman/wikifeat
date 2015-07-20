#!/bin/sh
# setup_users.sh
# User setup script
# 
# Sets up the user database in couchdb 
# Format: setup_users.sh hostname portnumber

if [ $# -eq 0 ]
then
	echo
	echo "Sets up couchdb user database"
	echo "specify the hostname and port of your couhdb server."
	echo "format is setup_users.sh hostname port."
	echo
	echo "Assuming localhost:5984"
	echo
fi

host=${1-localhost}
port=${2-5984}
admin_user=$3
admin_password=$4

if [ -z "$3" ]
then

	echo -n "Enter couchdb admin username: "
	read admin_user
fi
if [ -z "$4" ]
then
	echo -n "Enter couchdb admin password: " 
	stty -echo
	read admin_password
	stty echo
fi
echo

#Set the user public fields
echo -n "Setting user database public fields"
url=http://$host:$port/_config/couch_httpd_auth/public_fields
content_type="Content-Type: application/json"
field=userPublic
auth=$admin_user:$admin_password
public_user_fields_cmd="curl -X PUT $url -H \"$content_type\" -d \"$field\" -u $auth"
$public_user_fields_cmd

#Set up views in the user database
echo -n "Creating design document for user queries"
url=http://$host:$port/_users/_design/user_queries
echo -n $userDdoc
writeDesignDocCmd="curl -X PUT $url -H \"$content_type\" --data-binary @ddoc/user_ddoc.json -u $auth"
echo -n $writeDesignDocCmd
$writeDesignDocCmd

