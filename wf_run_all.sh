#!/bin/sh

ORGPATH=$GOPATH/src/github.com/rhinoman

cd $ORGPATH/wikifeat-users
./wikifeat-users&
cd $ORGPATH/wikifeat-wikis
./wikifeat-wikis&
cd $ORGPATH/wikifeat-frontend
./wikifeat-frontend&
