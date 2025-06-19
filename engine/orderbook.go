package engine

import (
	"TradeIT/models"
	"container/heap"
	"errors"
	"fmt"
	"sync"
)

type Orderbook struct {
	buy_orders  map[float64]*DoublyLinkedList
	sell_orders map[float64]*DoublyLinkedList
	bids_prices MaxHeap
	asks_prices MinHeap
	tradeCount  int64
	buyCount    int64
	sellcount   int64
	mu          sync.RWMutex
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func InitOrderBook_() *Orderbook {
	ob := Orderbook{
		buy_orders:  make(map[float64]*DoublyLinkedList),
		sell_orders: make(map[float64]*DoublyLinkedList),
	}
	heap.Init(&ob.bids_prices)
	heap.Init(&ob.asks_prices)
	return &ob
}

func (ob *Orderbook) InsertOrder(orderData models.Metadata) error {
	ob.mu.Unlock()
	defer ob.mu.Lock()
	if orderData.Order_type == "buy" {
		if _, exists := ob.buy_orders[orderData.Price]; !exists {
			ob.buy_orders[orderData.Price] = &DoublyLinkedList{}
			heap.Push(&ob.bids_prices, orderData.Price) // inserts the price in Heap
		}
		ob.buy_orders[orderData.Price].PushBack(orderData) //insert order details in buy order list
		return nil
	} else if orderData.Order_type == "sell" {
		if _, exists := ob.sell_orders[orderData.Price]; !exists {
			ob.sell_orders[orderData.Price] = &DoublyLinkedList{}
			heap.Push(&ob.asks_prices, orderData.Price) // inserts the price in Heap
		}
		ob.sell_orders[orderData.Price].PushBack(orderData) //insert order details in buy order list
		return nil
	} else {
		return errors.New("invalid order type")
	}
}

func (ob *Orderbook) Matcher(order models.Metadata) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	if order.Order_type == "buy" { //buy order
		for ob.asks_prices.Len() > 0 && order.Remq > 0 {
			bestAsk := ob.asks_prices.Peek()
			if order.Price < bestAsk { //buyers price is
				break
			}

			askList := ob.sell_orders[bestAsk]
			for node := askList.Head; node != nil && order.Remq > 0; node = node.Next {
				matchQuantity := min(node.Metadata.Remq, order.Remq)
				order.Remq -= matchQuantity
				node.Metadata.Remq -= matchQuantity
				ob.tradeCount++
				// register the trade below
				//
				if node.Metadata.Remq == 0 {
					if node.Metadata.Remq == 0 {
						// Update linked list pointers
						if node.Prev != nil {
							node.Prev.Next = node.Next
						} else {
							askList.Head = node.Next
						}

						if node.Next != nil {
							node.Next.Prev = node.Prev
						} else {
							askList.Tail = node.Prev
						}

						askList.Size--
						ob.sellcount--
					}
				}

				if order.Remq == 0 {
					continue
				}

			}
			if askList.Size == 0 {
				delete(ob.sell_orders, bestAsk)
				heap.Pop(&ob.asks_prices)
			}
		}
		if order.Remq > 0 { //adds the order to the orderbook if not filled
			ob.InsertOrder(order)
			ob.buyCount++
		}
	} else {
		// Proccessing the sell order

		for ob.bids_prices.Len() > 0 && order.Remq > 0 { //checking for the bids available
			bestBid := ob.bids_prices.Peek()
			if order.Price > bestBid { //breaks if bids are less than the asks
				break
			}
			bidList := ob.buy_orders[bestBid]
			node := bidList.Head
			for node != nil && order.Remq > 0 {
				//trade occurs
				matchQuantity := min(order.Remq, node.Metadata.Remq)
				order.Remq -= matchQuantity
				node.Metadata.Remq -= matchQuantity
				ob.tradeCount++
				//register the trade here
				//
				if node.Metadata.Remq == 0 { //buyer order is fullfilled and node has become stale
					if node.Metadata.Remq == 0 {
						// Update linked list pointers
						if node.Prev != nil {
							node.Prev.Next = node.Next
						} else {
							bidList.Head = node.Next
						}

						if node.Next != nil {
							node.Next.Prev = node.Prev
						} else {
							bidList.Tail = node.Prev
						}

						bidList.Size--
						ob.buyCount--
					}
				}

				if order.Remq == 0 {
					continue
				}
				node = node.Next
			}
			if bidList.Size == 0 {
				delete(ob.buy_orders, bestBid)
				heap.Pop(&ob.bids_prices)
			}
		}
		if order.Remq > 0 {
			ob.InsertOrder(order)
			ob.sellcount++
		}

	}

	if ob.asks_prices.Len() == 0 && ob.bids_prices.Len() == 0 {
		//close the matcher and orderbook
		fmt.Println("Time to close the orderbook and matcher")
	}

}

func (ob *Orderbook) DisplayResult() {
	fmt.Printf("Trades: %d |\nBuyOrders: %d|\nSellOrders: %d|\nTotal: %d", ob.tradeCount, ob.buyCount, ob.sellcount, (ob.tradeCount + ob.sellcount + ob.buyCount))
	fmt.Println("\nHeap lengths: ")
	fmt.Printf("bids: %d\nasks: %d", ob.bids_prices.Len(), ob.asks_prices.Len())
	fmt.Println("")
}
