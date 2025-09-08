package engine

import (
	"TradeIT/models"
	"fmt"
	"log"
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

func (ld *Ledger) ProcessOrder(order models.Order) error {
	ob, exists := ld.ent[order.Stock]
	fmt.Println(order.Id, " reached the ledger")
	if !exists {
		ob = InitOrderBook()
		ld.ent[order.Stock] = ob
		ld.ob_count++
	}
	err := ob.Matcher(order)
	if err != nil {
		log.Println("error occured in Matcher ", err)
	}
	return err
}

func (ld *Ledger) CancelOrder(stock string, order_id uint64) error {
	ob, exists := ld.ent[stock]
	if !exists {
		log.Println("Order not found")
	}
	return ob.CancelOrder(order_id)
}
