#!/bin/sh

WIKIPATH=$GOPATH/src/github.com/rhinoman/wikifeat

cd $WIKIPATH/auth
./auth&
sleep 2
cd $WIKIPATH/users
./users&
sleep 2
cd $WIKIPATH/wikis
./wikis&
sleep 2
cd $WIKIPATH/notifications
./notifications&
sleep 2
cd $WIKIPATH/frontend
./frontend&
