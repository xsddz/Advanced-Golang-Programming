package main

import (
	"gin-app/library/app"
	"gin-app/router"
)

func main() {
	a := app.Init()
	a.RegisterServer(app.HTTPServerMaker(router.SetGinRouter))
	a.RegisterServer(app.RCPServerMaker(router.SetGRPCRouter))
	a.Run()
}
