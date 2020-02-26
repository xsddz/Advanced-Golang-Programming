
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wsserver_linux64 main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o wsserver_mac64 main.go
