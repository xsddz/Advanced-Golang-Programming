package main

import (
	"gin-app/library/server"
	"gin-app/router"
)

func main() {
	app := server.NewAPP()
	app.RegisterServer(server.NewHTTPServer(router.HTTPRouter))
	app.RegisterServer(server.NewGRPCServer(router.GRPCRouter))
	app.Run()
}
