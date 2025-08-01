package engine

import (
	"TradeIT/models"
	"sync"
)

type Ledger struct {
	ent sync.Map
	mu  sync.Mutex
}

func InitLedger() *Ledger{
	return &Ledger{}
}

func (ld *Ledger) InitOrderBook(stock string) {
	ld.ent.Store(stock,InitOrderBook_())        
}

func (ld* Ledger) ProcessOrder(order models.Metadata){
	
}
