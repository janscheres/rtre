package main

import (
	"log"

)

func main() {
	log.Println("Initialisng Real-Time Risk Engine")

	server := riskServer{}

	go startgRPCServer(&server)
	select {
	}
}
