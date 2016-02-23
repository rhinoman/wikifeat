#!/bin/sh
# build_all.sh
# Build everything and copy it into a build folder

# Because cgo hates clang.  This is necessary to build on *nixes 
# that don't use GNU stuff by default, like FreeBSD
CC=gcc
export CC

VERSION=0.7.0-alpha
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
mkdir ${BUILD_DIR}/auth
mkdir ${BUILD_DIR}/config
mkdir ${BUILD_DIR}/users
mkdir ${BUILD_DIR}/scripts
mkdir ${BUILD_DIR}/scripts/db_update
mkdir ${BUILD_DIR}/wikis
mkdir ${BUILD_DIR}/notifications
mkdir ${BUILD_DIR}/frontend
mkdir ${BUILD_DIR}/frontend/web_app
# Build stuff
go build -v -o ${BUILD_DIR}/config/wikifeat-config ../config
go build -v -o ${BUILD_DIR}/auth/wikifeat-auth ../auth
go build -v -o ${BUILD_DIR}/users/wikifeat-users ../users
go build -v -o ${BUILD_DIR}/wikis/wikifeat-wikis ../wikis
go build -v -o ${BUILD_DIR}/notifications/wikifeat-notifications ../notifications
go build -v -o ${BUILD_DIR}/frontend/wikifeat-frontend ../frontend
# Copy some supporting files
cp ../config/config.ini.example ${BUILD_DIR}/config/config.ini.example
cp ../frontend/index.html.template ${BUILD_DIR}/frontend/index.html.template
cp -R ../frontend/plugins ${BUILD_DIR}/frontend/plugins
rm -f ${BUILD_DIR}/frontend/plugins/plugins.ini
cp -R ../notifications/templates ${BUILD_DIR}/notifications/templates
# "Compile" the webapp, then copy it to the build dir
r.js -o ../frontend/web_app/app/scripts/app.build.js
cp -R ../frontend/web_app/wikifeat-build ${BUILD_DIR}/frontend/web_app/app
rm -rf ../frontend/web_app/wikifeat-build
# Copy some scripts, man
cp -p ../wf_run_all.sh ${BUILD_DIR}/wf_run_all.sh
cp -p ../kill_all.sh ${BUILD_DIR}/kill_all.sh
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

