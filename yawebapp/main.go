package main

import (
	"yawebapp/library/inner/app"
	"yawebapp/library/inner/server"
	"yawebapp/routers"
)

func main() {
	app.Init()
	app.RegisterServer(server.NewHTTPServer(routers.HTTPRouter))
	app.RegisterServer(server.NewGRPCServer(routers.GRPCRouter))
	app.Run()
}
