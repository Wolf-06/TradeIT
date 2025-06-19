package engine

import (
	"TradeIT/models"
	"container/heap"
	"sync"
)

type Ledger struct {
	ent map[string]*Orderbook
	mu  sync.Mutex
}

func (ld *Ledger) InitOrderBook(key string) {
	ld.mu.Lock()
	defer ld.mu.Unlock()
	ld.ent[key] = &Orderbook{
		buy_orders:  make(map[float64]*DoublyLinkedList),
		sell_orders: make(map[float64]*DoublyLinkedList),
	}
	heap.Init(&ld.ent[key].asks_prices) //initializes the heaps
	heap.Init(&ld.ent[key].bids_prices)
}

func (ld *Ledger) InsertEntry(orderData models.Metadata) {
	if _, exists := ld.ent[orderData.Stock]; !exists {
		ld.InitOrderBook(orderData.Stock)
	}

	ld.ent[orderData.Stock].Matcher(orderData)
	//if err != nil {
	//	fmt.Println("Error in inserting order in OrderBook: ", err)
	//}
}
