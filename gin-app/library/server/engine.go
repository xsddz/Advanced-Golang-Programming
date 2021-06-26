package server

import (
	"fmt"
	"sync"
)

type ServerHandle func(*sync.WaitGroup)

type Engine struct {
	servers []ServerHandle
}

func Init() *Engine {
	return &Engine{}
}

func (e *Engine) RegisterServer(s ServerHandle) {
	e.servers = append(e.servers, s)
}

func (e *Engine) Run() {
	var wg sync.WaitGroup
	wg.Add(len(e.servers))

	for _, sh := range e.servers {
		go sh(&wg)
	}

	wg.Wait()
	fmt.Println("the end.")
}
