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
# Build stuff
go build -o ../build/users/users ../users
go build -o ../build/wikis/wikis ../wikis 
go build -o ../build/frontend/frontend ../frontend 
# Copy some supporting files
cp ../users/config.ini ../build/users/config.ini
cp ../wikis/config.ini ../build/wikis/config.ini
cp ../frontend/config.ini ../build/frontend/config.ini
# Copy the web app for the frontend service
cp -R ../frontend/web_app ../build/frontend/web_app
