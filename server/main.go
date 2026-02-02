package main

import (
	"log"

)

func main() {
	wsclient := WsClient{
		messages: make(chan []byte, 100),
		orderbook: OrderBook{
			Bids: make(map[float64]float64),
			Asks: make(map[float64]float64),
		},
	}

	go startgRPCServer(&wsclient.orderbook)

	for {
		wsclient.connect()

		log.Println("[NET] Connection died, restarting...")
	}
}
