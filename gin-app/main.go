package main

import (
	"fmt"
	"gin-app/router"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go StartRPCServer(&wg)
	go StartHTTPServer(&wg)

	wg.Wait()
	fmt.Println("the end.")
}

func StartRPCServer(wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("rpc server end.")
		wg.Done()
	}()

}

func StartHTTPServer(wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("http server end.")
		wg.Done()
	}()

	app := gin.Default()

	router.SetGinRouter(app)

	// app.Run()
	fmt.Printf("Listening and serving HTTP on %s\n", ":8080")
	err := http.ListenAndServe(":8080", app)
	if err != nil {
		fmt.Println("Start HTTP server error:", err)
	}
}
