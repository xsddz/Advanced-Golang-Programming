package main

import (
	"gin-app/library/app"
	"gin-app/library/server"
	"gin-app/router"
)

func main() {
	app.Init([]string{"mysql", "redis"})
	app.RegisterServer(server.NewHTTPServer(router.HTTPRouter))
	app.RegisterServer(server.NewGRPCServer(router.GRPCRouter))
	app.Run()
}
