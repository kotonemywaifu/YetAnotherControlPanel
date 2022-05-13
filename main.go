package main

import (
	"log"

	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/server"
)

func main() {
	log.Println("Launching YetAnotherControlPanel...")

	log.Println("Initializing config...")
	err := others.InitConfig()
	if err != nil {
		panic(err)
	}

	log.Println("Starting server...")
	server.StartServer()

	// keep goroutine alive
	select {}
}