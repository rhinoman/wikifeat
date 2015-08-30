#!/bin/sh
# build_all.sh
# Build everything and copy it into a build folder

# Because cgo hates clang.  This is necessary to build on *nixes 
# that don't use GNU stuff by default, like FreeBSD
CC=gcc
export CC

VERSION=0.1-alpha
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
mkdir ${BUILD_DIR}/wikis
mkdir ${BUILD_DIR}/frontend
mkdir ${BUILD_DIR}/frontend/web_app
# Build stuff
go build -v -o ${BUILD_DIR}/users/users ../users
go build -v -o ${BUILD_DIR}/wikis/wikis ../wikis 
go build -v -o ${BUILD_DIR}/frontend/frontend ../frontend 
# Copy some supporting files
cp ../users/config.ini.example ${BUILD_DIR}/users/config.ini.example
cp ../wikis/config.ini.example ${BUILD_DIR}/wikis/config.ini.example
cp ../frontend/config.ini.example ${BUILD_DIR}/frontend/config.ini.example
cp -R ../frontend/plugins ${BUILD_DIR}/frontend/plugins
# "Compile" the webapp, then copy it to the build dir
r.js -o ../frontend/web_app/app/scripts/app.build.js
cp -R ../frontend/web_app/wikifeat-build ${BUILD_DIR}/frontend/web_app/app
rm -rf ../frontend/web_app/wikifeat-build
# Copy some scripts, man
cp -p ../wf_run_all.sh ${BUILD_DIR}/wf_run_all.sh
cp -R ../scripts/ddoc ${BUILD_DIR}/scripts
cp -p ../scripts/setup_users.sh ${BUILD_DIR}/scripts
cp -p ../scripts/create_master_user.sh ${BUILD_DIR}/scripts
cp -p ../scripts/setup.sh ${BUILD_DIR}/scripts
cp -p ../scripts/get_rev.sh ${BUILD_DIR}/scripts
# Now make a tarball
mkdir ../dist
echo "Creating ${TARNAME}"
cd ../build
tar cvzf ../dist/${TARNAME} ${BUILDNAME}
echo "All Done!"

