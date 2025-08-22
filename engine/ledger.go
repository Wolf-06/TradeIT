package engine

import (
	"TradeIT/models"
)

type Ledger struct {
	ent      map[string]*Orderbook
	ob_count int64
}

func InitLedger() *Ledger {
	return &Ledger{
		ent: make(map[string]*Orderbook),
	}
}

func (ld *Ledger) ProcessOrder(order models.Metadata) error {
	ob, exists := ld.ent[order.Stock]
	if !exists {
		ob = InitOrderBook()
		ld.ent[order.Stock] = ob
		ld.ob_count++
	}
	err := ob.Matcher(order)
	return err
}
