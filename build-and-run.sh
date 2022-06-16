#!/bin/bash
set -exu

cd ./src
go build -a -v -o ../build .
cd ../build
sudo ./netlink-api
