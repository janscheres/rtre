package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	websocket *websocket.Conn

	messages chan []byte
	
	done chan struct{}
}

var wsclient WsClient

func receive() {
	defer close(wsclient.done)
	for {
		_, msg, err := wsclient.websocket.ReadMessage()
		if err != nil {
			log.Println("[RCV]: Error recieving message", err)
			return
		}

		wsclient.messages <- msg
	}
}

func parseAndPass() {
	defer close(wsclient.done)
	for {
		msg := <- wsclient.messages

		var update Answer
		err := json.Unmarshal(msg, &update)
		if err != nil {
			log.Fatal("[JSON] Error parsing json", err)
		}

		log.Println(update.Symbol, update.Bids[0])

	}

}

func main() {
	ws, res, err := websocket.DefaultDialer.Dial("wss://fstream.binance.com/ws/btcusdt@depth@100ms", nil)
	if err != nil {
		log.Fatal("[DIAL] Couldn't connect to Binance API:", err)
	}
	wsclient.websocket = ws
	wsclient.messages = make(chan []byte)
	defer wsclient.websocket.Close()

	wsclient.done = make(chan struct{})

	go receive()

	go parseAndPass()

	<-wsclient.done

	log.Println(ws, res, err)
}
