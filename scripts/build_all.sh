#!/bin/sh
# build_all.sh
# Build everything and copy it into a build folder

# Because cgo hates clang.  This is necessary to build on *nixes 
# that don't use GNU stuff by default, like FreeBSD
CC=gcc
export CC

VERSION=0.3.1-alpha
ARCH=`uname -p`
OS=`uname`
BUILDNAME=wikifeat_${VERSION}.${OS}-${ARCH}
TARNAME=${BUILDNAME}.tar.gz
BUILD_DIR=../build/${BUILDNAME}
# erase any previous build
rm -rf ../build

# Make some directories
mkdir ../build
mkdir $BUILD_DIR
mkdir ${BUILD_DIR}/users
mkdir ${BUILD_DIR}/scripts
mkdir ${BUILD_DIR}/scripts/db_update
mkdir ${BUILD_DIR}/wikis
mkdir ${BUILD_DIR}/notifications
mkdir ${BUILD_DIR}/frontend
mkdir ${BUILD_DIR}/frontend/web_app
# Build stuff
go build -v -o ${BUILD_DIR}/users/users ../users
go build -v -o ${BUILD_DIR}/wikis/wikis ../wikis 
go build -v -o ${BUILD_DIR}/notifications/notifications ../notifications
go build -v -o ${BUILD_DIR}/frontend/frontend ../frontend 
# Copy some supporting files
cp ../users/config.ini.example ${BUILD_DIR}/users/config.ini.example
cp ../wikis/config.ini.example ${BUILD_DIR}/wikis/config.ini.example
cp ../notifications/config.ini.example ${BUILD_DIR}/notifications/config.ini.example
cp ../frontend/config.ini.example ${BUILD_DIR}/frontend/config.ini.example
cp -R ../frontend/plugins ${BUILD_DIR}/frontend/plugins
cp -R ../notifications/templates ${BUILD_DIR}/notifications/templates
# "Compile" the webapp, then copy it to the build dir
r.js -o ../frontend/web_app/app/scripts/app.build.js
cp -R ../frontend/web_app/wikifeat-build ${BUILD_DIR}/frontend/web_app/app
rm -rf ../frontend/web_app/wikifeat-build
# Copy some scripts, man
cp -p ../wf_run_all.sh ${BUILD_DIR}/wf_run_all.sh
cp -R ../scripts/ddoc ${BUILD_DIR}/scripts
cp -R ../scripts/libs ${BUILD_DIR}/scripts
cp -p ../scripts/setup.py ${BUILD_DIR}/scripts
cp -p ../scripts/util.py ${BUILD_DIR}/scripts
cp -p ../scripts/config.py ${BUILD_DIR}/scripts
cp -p ../scripts/install.py ${BUILD_DIR}/scripts
cp ../scripts/db_update/*.py ${BUILD_DIR}/scripts/db_update
cp ../scripts/db_update/*.md ${BUILD_DIR}/scripts/db_update
# Now make a tarball
mkdir ../dist
echo "Creating ${TARNAME}"
cd ../build
tar cvzf ../dist/${TARNAME} ${BUILDNAME}
echo "All Done!"

