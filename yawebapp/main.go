package main

import (
	"yawebapp/library/common/app"
	"yawebapp/library/common/server"
	"yawebapp/routers"
)

func main() {
	app.Init()
	app.RegisterServer(server.NewHTTPServer(routers.HTTPRouter))
	app.RegisterServer(server.NewGRPCServer(routers.GRPCRouter))
	app.Run()
}
