package main

import (
	"log"

	//pb "github.com/janscheres/rtre/pb"
)

func main() {
	wsclient := WsClient{
		messages: make(chan []byte, 100),
	}

	for {
		wsclient.connect()

		log.Println("[NET] Connection died, restarting...")
	}
}
