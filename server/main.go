package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

type Order struct {
	Price float64
	Quantity float64
}

func (o *Order) UnmarshalJSON(data []byte) error {
	var raw [2]string
	err := json.Unmarshal(data, &raw)
	if err != nil {
		log.Println("ERROR: [JSON] Error parsing order", err)
	}

	o.Price, err = strconv.ParseFloat(raw[0], 64)
	if err != nil {
		log.Println("ERROR: Converting price to float", err)
	}
	o.Quantity, err = strconv.ParseFloat(raw[1], 64)
	if err != nil {
		log.Println("ERROR: Converting quantity to int", err)
	}

	return nil
}

type Answer struct {
	EventType string		`json:"e"`
	EventTime int64			`json:"E"`
	TransactionTime int64	`json:"T"`
	Symbol string			`json:"s"`
	FirstUpdateID int64		`json:"U"`
	LastUpdateID int64		`json:"u"`
	PrevFinalID int64		`json:"pu"`
	Bids []Order			`json:"b"`
	Asks []Order			`json:"a"`
}

type WsClient struct {
	websocket *websocket.Conn

	messages chan []byte
	
	done chan struct{}
}

func (c *WsClient) receive() {
	defer close(c.done)
	for {
		_, msg, err := c.websocket.ReadMessage()
		if err != nil {
			log.Println("ERROR: [RCV]: Error recieving message", err)
			close(c.messages)
			return
		}

		c.messages <- msg
	}
}

func (c *WsClient) parseAndPass() {
	defer close(c.done)
	for msg := range c.messages {
		var update Answer

		err := json.Unmarshal(msg, &update)
		if err != nil {
			log.Println("ERROR: [JSON] Error parsing json", err)
		}

		//log.Println(update.Symbol, update.Bids[0])
	}

}

func main() {
	var wsclient WsClient

	ws, res, err := websocket.DefaultDialer.Dial("wss://fstream.binance.com/ws/btcusdt@depth@100ms", nil)
	if err != nil {
		log.Fatal("[DIAL] Couldn't connect to Binance API:", err)
	}
	wsclient.websocket = ws
	wsclient.messages = make(chan []byte, 100)
	defer wsclient.websocket.Close()

	wsclient.done = make(chan struct{})

	go wsclient.receive()

	go wsclient.parseAndPass()

	<-wsclient.done

	log.Println(ws, res, err)
}
