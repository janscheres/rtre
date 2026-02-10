package main

import (
	"log"
	"sync"
)

type OrderBook struct {
	mu sync.RWMutex

	Bids map[float64]float64
	Asks map[float64]float64

	TotalBidVol float64
	TotalAskVol float64

	OBIChan chan float64
	SpreadChan chan float64
}

func (b *OrderBook) handleUpdate(update Answer) {
	b.mu.Lock()

	var maxBid float64
	for _, u := range update.Bids {
		b.TotalBidVol-=b.Bids[u.Price]

		if u.Quantity == 0 {
			delete(b.Bids, u.Price)
		} else if u.Quantity > 0 {
			if u.Price > maxBid {
				maxBid = u.Price
			}

			b.Bids[u.Price]=u.Quantity
			b.TotalBidVol+=u.Quantity
		}
		
	}

	var maxAsk float64
	for _, u := range update.Asks {
		b.TotalAskVol-=b.Asks[u.Price]

		if u.Quantity == 0 {
			delete(b.Asks, u.Price)
		} else if u.Quantity > 0 {
			if u.Price > maxAsk {
				maxAsk = u.Price
			}

			b.Asks[u.Price]=u.Quantity
			b.TotalAskVol+=u.Quantity
		}
	}

	b.mu.Unlock()

	select {
	case b.SpreadChan <- maxAsk-maxBid:
	default:
		log.Println("[go] DROPPING spread due to full channel :(")
	}

	select {
	case b.OBIChan <- b.GetOBI():
	default:
		log.Println("[go] DROPPING obi due to full channel :(")
	}
}

func (b *OrderBook) GetOBI() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	bidVol, askVol := b.TotalBidVol, b.TotalAskVol

	total := bidVol+askVol

	if total == 0 {
		return 0// return 0 for balanced since we have no orders
	}
	return (bidVol-askVol)/total
}

