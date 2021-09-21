package main

import (
	"yawebapp/library/infra/app"
	"yawebapp/routers"
)

func main() {
	app.Init()
	app.AddServerRouter("http", routers.HTTPRouter)
	app.AddServerRouter("grpc", routers.GRPCRouter)
	app.Run()
}
