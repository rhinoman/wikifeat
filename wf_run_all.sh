#!/bin/sh

WIKIPATH=$GOPATH/src/github.com/rhinoman/wikifeat

cd $WIKIPATH/config
./wikifeat-config&
cd $WIKIPATH/auth
./wikifeat-auth&
sleep 2
cd $WIKIPATH/users
./wikifeat-users&
sleep 2
cd $WIKIPATH/wikis
./wikifeat-wikis&
sleep 2
cd $WIKIPATH/notifications
./wikifeat-notifications&
sleep 2
cd $WIKIPATH/frontend
./wikifeat-frontend&
