#!/bin/bash

## 1. brew install protobuf
## 2. go install github.com/golang/protobuf/protoc-gen-go
## 3. export PATH="$PATH:$(go env GOPATH)/bin"

OLD_DIR=`pwd`
CUR_DIR=`dirname $0`
cd $CUR_DIR

protoc --go_out=plugins=grpc:../ pb_protos/*.proto

cd $OLD_DIR
