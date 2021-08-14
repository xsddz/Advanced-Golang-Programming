package main

import (
	"yawebapp/Illuminate/app"
	"yawebapp/Illuminate/server"
	"yawebapp/router"
)

func main() {
	app.Init()
	app.RegisterServer(server.NewHTTPServer(router.HTTPRouter))
	app.RegisterServer(server.NewGRPCServer(router.GRPCRouter))
	app.Run()
}
