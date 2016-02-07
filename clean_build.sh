#!/bin/sh

cd auth
go clean; go build -o wikifeat-auth
cd ..
cd wikis
go clean; go build -o wikifeat-wikis
cd ..
cd users
go clean; go build -o wikifeat-users
cd ..
cd notifications
go clean; go build -o wikifeat-notifications
cd ..
cd frontend
go clean; go build -o wikifeat-frontend
cd ..
