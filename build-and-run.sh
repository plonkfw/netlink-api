#!/bin/bash
set -exu

cd ./src
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -o ../build .
cd ../build
file ./netlink-api | tr , '\n'

# defaults to localhost:4821 if not provided
sudo LISTEN=:4821 ./netlink-api
