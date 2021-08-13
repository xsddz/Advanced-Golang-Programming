package app

import (
	"fmt"
	"gin-app/library/server"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/shima-park/agollo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	defaultApp *server.Engine

	APOLLO agollo.Agollo

	SQLITE *gorm.DB
	REDIS  *redis.Client
	DB     *gorm.DB
)

func Init(keys []string) {
	initApollo()

	for _, key := range keys {
		if init, ok := initTable[strings.ToLower(key)]; ok {
			fmt.Println("init: ", key)
			init()
		}
	}

	defaultApp = server.NewEngine(APOLLO.Get("APP_ENV"))
}

func RegisterServer(s server.ServerI) {
	defaultApp.RegisterServer(s)
}

func GetHTTPServer() *gin.Engine {
	return defaultApp.GetHTTPServer()
}

func GetGRPCServer() *grpc.Server {
	return defaultApp.GetGRPCServer()
}

func Run() {
	defaultApp.Run()
}
