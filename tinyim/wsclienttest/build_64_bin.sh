
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wsclient_linux64 main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o wsclient_mac64 main.go
