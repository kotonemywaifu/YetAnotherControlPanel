package task

import (
	"log"
	"time"
)

func StartTicking() {
	for {
		time.Sleep(time.Second * 1)
		go tick()
	}
}

func tick() {
	// don't panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error thrown in tick:", r)
		}
	}()

	for _, v := range tasks {
		v.Tick()
	}
}
