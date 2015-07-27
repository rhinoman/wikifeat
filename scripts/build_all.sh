#!/bin/sh
# build_all.sh
# Build everything and copy it into a build folder

# Because cgo hates clang.  This is necessary to build on *nixes 
# that don't use GNU stuff by default, like FreeBSD
CC=gcc
export CC

# Make some directories
mkdir ../build
mkdir ../build/users
mkdir ../build/wikis
mkdir ../build/frontend
mkdir ../build/frontend/web_app
# Build stuff
go build -v -o ../build/users/users ../users
go build -v -o ../build/wikis/wikis ../wikis 
go build -v -o ../build/frontend/frontend ../frontend 
# Copy some supporting files
cp ../users/config.ini ../build/users/config.ini
cp ../wikis/config.ini ../build/wikis/config.ini
cp ../frontend/config.ini ../build/frontend/config.ini
cp -R ../frontend/plugins ../build/frontend/plugins
# Copy the web app for the frontend service
r.js -o ../frontend/web_app/app/scripts/app.build.js
cp -R ../frontend/web_app/wikifeat-build ../build/frontend/web_app/app
rm -rf ../frontend/web_app/wikifeat-build
