package main

import (
	"log"

	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/server"
)

func main() {
	log.Println("Initializing config...")
	err := others.InitEnv()
	if err != nil {
		log.Panicln(err)
	}
	err = others.InitConfig()
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Starting server...")
	server.StartServer()

	// keep goroutine alive
	select {}
}
