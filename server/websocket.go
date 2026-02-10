package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

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
	EventType 		string	`json:"e"`
	EventTime 		int64	`json:"E"`
	TransactionTime int64	`json:"T"`
	Symbol 			string	`json:"s"`
	FirstUpdateID 	int64	`json:"U"`
	LastUpdateID 	int64	`json:"u"`
	PrevFinalID 	int64	`json:"pu"`
	Bids 			[]Order	`json:"b"`
	Asks 			[]Order	`json:"a"`
}

type WsClient struct {
	ctx context.Context

	websocket *websocket.Conn
	orderbook OrderBook
	symbol string

	messages chan []byte
	done chan struct{}
}

func (c *WsClient) run() {
	for {
		select {
		case <-c.ctx.Done():
			log.Println("[WS] Client disconnected, stopping")
			return
		default:
			err := c.connect()
			if err != nil {
				log.Println("[NET] attempting reconnecting to upstream, was disconnected:", err)
				select {
				case <-time.After(time.Second):
				case <-c.ctx.Done():
					return
				}
			}
		}
	}
}

func (c *WsClient) receive() {
	defer close(c.done)
	defer close(c.messages)

	for {
		_, msg, err := c.websocket.ReadMessage()
		if err != nil {
			log.Println("ERROR: [RCV]: Error recieving message", err)
			return
		}


		select {
		case c.messages <-msg:
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *WsClient) parseAndPass() {
	for {
		select {
		case msg, ok := <-c.messages:
			if !ok {
				return
			}

			var update Answer
			err := json.Unmarshal(msg, &update)
			if err != nil {
				log.Println("ERROR: [JSON] Error parsing json", err)
			}

			c.orderbook.handleUpdate(update)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *WsClient) connect() error {
	c.messages = make(chan []byte, 100)
	c.done = make(chan struct{})

	ws, _, err := websocket.DefaultDialer.Dial("wss://fstream.binance.com/ws/"+c.symbol+"@depth@100ms", nil)
	if err != nil {
		//log.Println("[DIAL] Couldn't connect to Binance API:", err)
		return err
	}
	c.websocket = ws

	log.Println("Successfully connected to upstream websocket")
	
	go c.receive()
	go c.parseAndPass()

	select {
	case <-c.done:
	case <-c.ctx.Done():
		ws.Close()
	}

	return nil
}


