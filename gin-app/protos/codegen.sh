

## 1. brew install protobuf
## 2. go install github.com/golang/protobuf/protoc-gen-go
## 4. export PATH="$PATH:$(go env GOPATH)/bin"


protoc --go_out=plugins=grpc:../ $1

