package main

import (
	"log"

)

func main() {
	log.Println("Initialisng Real-Time Risk Engine")

	server := riskServer{}

	//go startgRPCServer(&wsclient.orderbook)
	go startgRPCServer(&server)
	for {
	}

	/*for {
		wsclient.connect()

		log.Println("[NET] Connection died, restarting...")
	}*/
}
