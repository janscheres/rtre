package main

import (
	"sync"
)

type OrderBook struct {
	mu sync.RWMutex

	Bids map[float64]float64
	Asks map[float64]float64

	TotalBidVol float64
	TotalAskVol float64

	OBIChan chan float64
}

func (b *OrderBook) handleUpdate(update Answer) {
	b.mu.Lock()

	for _, u := range update.Bids {
		b.TotalBidVol-=b.Bids[u.Price]

		if u.Quantity == 0 {
			delete(b.Bids, u.Price)
		} else if u.Quantity > 0 {

			b.Bids[u.Price]=u.Quantity
			b.TotalBidVol+=u.Quantity
		}
		
	}

	for _, u := range update.Asks {
		b.TotalAskVol-=b.Asks[u.Price]

		if u.Quantity == 0 {
			delete(b.Asks, u.Price)
		} else if u.Quantity > 0 {
			b.Asks[u.Price]=u.Quantity
			b.TotalAskVol+=u.Quantity
		}
	}

	b.mu.Unlock()

	b.OBIChan <- b.GetOBI()
}

func (b *OrderBook) GetOBI() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	bidVol, askVol := b.TotalBidVol, b.TotalAskVol

	total := bidVol+askVol

	if total == 0 {
		return -100
	}
	return (bidVol-askVol)/total
}

