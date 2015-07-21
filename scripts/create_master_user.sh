#!/bin/sh
# create_master_user.sh
# Creates a Master User for Wikifeat

set -e

if [ $# -eq 0 ]
then
	echo
	echo "Create a Master User for Wikifeat"
	echo "You should have one (preferably only one) master user"
fi

host=${1-localhost}
port=${2-5984}
admin_user=$3
admin_password=$4
main_db=$5

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
if [ -z "$5" ]
then 
	echo -n "Enter the name of the main wikifeat database: "
	read main_db
fi

couch_host=http://$host:$port
auth=$admin_user:$admin_password
content_type="Content-Type: application/json"

#Create master user
while [ 1 ]
do 
	echo
	echo "You need at least (preferably *only*) one master user"
	echo
	echo -n "Enter master username: " 
	read master_username 
	echo
	echo -n "Enter Last Name: "
	read last_name
	echo
	echo -n "Enter First Name: "
	read first_name
	echo
	echo -n "Enter password for $master_username: " 
	stty -echo 
	read master_password 
	echo
	stty echo 
	echo -n "Verify password: " 
	stty -echo 
	read verify_password 
	stty echo 
	if [ $master_password != $verify_password ] 
	then 
		echo
		echo "Passwords don't match!"
	       	echo "Try again" 
		echo
	else 
		break
	fi
done
echo

url=http://$host:$port/_users/org.couchdb.user:$master_username

user_data="{\
  \"_id\":\"org.couchdb.user:$master_username\",\
  \"name\":\"$master_username\",\
  \"type\":\"user\",\
  \"roles\": [\"master\",\"admin\", \"$main_db:admin\"],\
  \"password\": \"$master_password\",\
  \"userPublic\": {\"lastName\":\"$last_name\",\"firstName\":\"$first_name\"}\
}"

echo -n "Creating the master user"
echo
curl -sS -X PUT $url -H "$content_type" -d "$user_data" -u $auth > /dev/null
echo
echo

