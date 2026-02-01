package main

import (
	"encoding/json"
	"log"
)

type Order struct {
	Price string
	Quantity string
}

func (o *Order) UnmarshalJSON(data []byte) error {
	var raw [2]string
	err := json.Unmarshal(data, &raw)
	if err != nil {
		log.Fatal("[JSON] Error parsing order", err)
	}

	o.Price = raw[0]
	o.Quantity = raw[1]
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
