package main

import (
	"log"

	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/server"
	"github.com/liulihaocai/YetAnotherControlPanel/task"
)

func main() {
	err := others.InitEnv()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Initializing config...")
	err = others.InitConfig()
	if err != nil {
		log.Panicln(err)
	}
	err = others.InitAccounts()
	if err != nil {
		log.Panicln(err)
	}
	err = others.SetupLogger()
	if err != nil {
		log.Panicln(err)
	}

	go task.StartTicking()

	log.Println("Starting server...")
	server.StartServer()
}
