#!/bin/sh
# setup.sh
# Master setup script
#
# Sets up Wikifeat prior to first run
# Format setup.sh hostname portnumber

set -e

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
avatar_db=user_avatars

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
echo "Creating Wikifeat Main Database $main_db"
echo
curl -sS -X PUT $couch_host/$main_db -u $auth 

# Load the main database design docs
echo "Loading Main Database Design Doc"
echo
url=$couch_host/$main_db/_design/wiki_query
echo "Check if Main database design doc already exists"
ddocRev=`sh ./get_rev.sh "$url" "$auth"`
if [ "$ddocRev" != "" ]; then
	revHeader="If-Match: $ddocRev"
	echo $revHeader
	echo "Design Doc Already exists: "$ddocRev"  Updating..."
	curl -sS -X PUT $url -H "$revHeader" -H "$content_type" --data-binary @ddoc/main_ddoc.json -u $auth 
else
	curl -sS -X PUT $url -H "$content_type" --data-binary @ddoc/main_ddoc.json -u $auth 
fi
# Load in the valdiation function
write_role="$main_db:write"
admin_role="admin"
validation_func="function(newDoc, oldDoc, userCtx){ \
if((userCtx.roles.indexOf(\\\"$write_role\\\") === -1) && \
(userCtx.roles.indexOf(\\\"$admin_role\\\") === -1) && \
(userCtx.roles.indexOf(\\\"master\\\") === -1) && \
(userCtx.roles.indexOf(\\\"_admin\\\") === -1)){ \
throw({forbidden: \\\"Not authorized\\\"}); \
}}" 
auth_doc="{ \
\"_id\": \"_design/_auth\", \
\"validate_doc_update\": \"$validation_func\" \
}"

url=$couch_host/$main_db/_security
echo "Setting security document for Main Database"
curl -sS -X PUT $url -H "$content_type" --data-binary @ddoc/main_access.json -u $auth

url=$couch_host/$main_db/_design/_auth
echo "Check if design doc exists"
ddocRev=`sh ./get_rev.sh "$url" "$auth"`
echo "Setting Write Access for Main Database"
if [ "$ddocRev" != "" ]; then
	revHeader="If-Match: $ddocRev"
	echo "Design Doc Already exists: $ddocRev  Updating..."
	curl -sS -X PUT $url -H "$revHeader" -H "$content_type" -d "$auth_doc" -u $auth 
else
	curl -sS -X PUT $url -H "$content_type" -d "$auth_doc" -u $auth 
fi
echo

#Create the user avatar database
echo "Creating User Avatar Database"
url=$couch_host/$avatar_db
curl -sS -X PUT $url -u $auth

url=$couch_host/$avatar_db/_security
echo "Setting security document for User Avatar Database"
curl -sS -X PUT $url -H "$content_type" --data-binary @ddoc/main_access.json -u $auth
echo "Setting Write access for User Avatar Database"
url=$couch_host/$avatar_db/_design/_auth
ddocRev=`sh ./get_rev.sh "$url" "$auth"`
if [ "$ddocRev" != "" ]; then
	revHeader="If-Match: $ddocRev"
	echo "Design Doc Already exists: $ddocRev  Updating..."
	curl -sS -X PUT $url -H "$revHeader" -H "$content_type" --data-binary @ddoc/avatar_auth.json -u $auth 
else
	curl -sS -X PUT $url -H "$content_type" --data-binary @ddoc/avatar_auth.json -u $auth 
fi


# Run the user db setup script
sh setup_users.sh $host $port $admin_user $admin_password $main_db

# Should we create a master user?

echo -n "Would you like to create a master user now (y/n)?"
read create_master
if echo "$create_master" | grep -iq "^y" ;then
	sh ./create_master_user.sh $host $port $admin_user $admin_password $main_db
fi

# TODO: Create config.ini ? 


echo
echo "Setup Complete!"


