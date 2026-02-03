package main

import (
	"log"

)

func main() {
	log.Println("Initialisng Real-Time Risk Engine")

	wsclient := WsClient{
		orderbook: OrderBook{
			Bids: make(map[float64]float64),
			Asks: make(map[float64]float64),
			OBIChan: make(chan float64, 100),
		},
	}

	go startgRPCServer(&wsclient.orderbook)

	for {
		wsclient.connect()

		log.Println("[NET] Connection died, restarting...")
	}
}
