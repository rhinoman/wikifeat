#!/bin/sh

WIKIPATH=$GOPATH/src/github.com/rhinoman/wikifeat

cd $WIKIPATH/users
./users&
sleep 2
cd $WIKIPATH/wikis
./wikis&
sleep 2
cd $WIKIPATH/frontend
./frontend&
