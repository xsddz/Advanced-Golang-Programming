package main

import (
	"gin-app/library/server"
	"gin-app/router"
)

func main() {
	app := server.Init()
	app.RegisterServer(server.HTTPServerMaker(router.SetGinRouter))
	app.RegisterServer(server.RCPServerMaker(router.SetGRPCRouter))
	app.Run()
}
