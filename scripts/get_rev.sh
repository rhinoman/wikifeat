#!/bin/sh

curl -Is $1 -u $2 | grep -Fi etag | awk '{gsub(/"/,"");gsub(/\r/,"");gsub(/ /,"");split($0,array,":")} END{print array[2]}'

