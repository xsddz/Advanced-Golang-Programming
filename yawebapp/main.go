package main

import (
	"yawebapp/library/infra/app"
	"yawebapp/routers"
)

func main() {
	app.Init()
	app.RegisterServer(app.NewHTTPServer(routers.HTTPRouter))
	app.RegisterServer(app.NewGRPCServer(routers.GRPCRouter))
	app.Run()
}
