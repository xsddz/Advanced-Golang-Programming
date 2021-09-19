#!/bin/bash
set -euo pipefail

## 安装依赖:
## 1. brew install protobuf
## 2. go install github.com/golang/protobuf/protoc-gen-go
## 3. export PATH="$PATH:$(go env GOPATH)/bin"

usage() {
    echo "Usage: protocoder.sh {workdir}"
    echo ""
}

if [ $# -ne 1 ]; then
    usage
    exit 1
fi

OLD_DIR=`pwd`
CUR_DIR=$1
cd $CUR_DIR

echo "process file:" *.proto
protoc --go_out=plugins=grpc:./ *.proto
echo "process end."

cd $OLD_DIR
