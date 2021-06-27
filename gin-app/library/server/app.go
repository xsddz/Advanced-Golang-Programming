package server

import (
	"fmt"
	"sync"
)

type ServerHandle func(*sync.WaitGroup)

type App struct {
	srvs []ServerHandle
}

func Init() *App {
	return &App{}
}

func (e *App) RegisterServer(s ServerHandle) {
	e.srvs = append(e.srvs, s)
}

func (e *App) Run() {
	var wg sync.WaitGroup
	wg.Add(len(e.srvs))

	for _, sh := range e.srvs {
		go sh(&wg)
	}

	wg.Wait()
	fmt.Println("the end.")
}
