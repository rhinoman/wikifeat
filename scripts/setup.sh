#!/bin/sh
# setup.sh
# Master setup script
#
# Sets up Wikifeat prior to first run
# Format setup.sh hostname portnumber

if [ $# -eq 0 ]
then
	echo
	echo "Sets up initial couchdb configuration."
	echo "specify the hostname and port of your couhdb server."
	echo "format is setup.sh hostname port."
	echo
	echo "Assuming localhost:5984"
	echo
fi

host=${1-localhost}
port=${2-5984}
admin_user=$3
admin_password=$4
main_db=wikifeat_main_db

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

while getopts  ":m:" opt; do
  case $opt in
	m) 
	  main_db=$OPTARG
	  ;;
	\?)
	  echo "Invalid option: -$OPTARG"
	  ;;
  esac
done

echo

couch_host=http://$host:$port
auth=$admin_user:$admin_password
content_type="Content-Type: application/json"

# Create the main database
echo -n "Creating Wikifeat Main Database $main_db"
curl -X PUT $couch_host/$main_db -u $auth

# Load the main database design docs
writeDesignDocCmd="curl -X PUT $couch_host/$main_db/_design/wiki_query -H \"$content_type\" --data-binary @ddoc/main_ddoc.json -u $auth"
$writeDesignDocCmd

# Load in the valdiation function
write_role="$main_db:write"
admin_role="$main_db:admin"
validation_func="function(newDoc, oldDoc, userCtx){
if((userCtx.roles.indexOf(\\\"$write_role\\\") === -1) && 
	(userCtx.roles.indexOf(\\\"$admin_role\\\") === -1) && 
	(userCtx.roles.indexOf(\\\"_admin\\\") === -1)){ 
	throw({forbidden: \\\"Not authorized\\\"}); 
}}" 
auth_doc="{ 
\"_id\": \"_design/_auth\", 
\"validate_doc_update\": \"$validation_func\" 
}"

echo -n $auth_doc
echo
echo $auth_doc > ddoc/main_security.json

writeAuthDocCmd="curl -X PUT $couch_host/$main_db/_design/_auth -H \"$content_type\" --data-binary @ddoc/main_security.json -u $auth"
echo -n $writeAuthDocCmd
echo
$writeAuthDocCmd
#`curl -X PUT $couch_host/$main_db/_design/_auth -H "$content_type" -d $auth_doc -u $auth`

echo

# Run the user db setup script
sh setup_users.sh $host $port $admin_user $admin_password

# TODO: Create config.ini ? 

