package engine

import (
	"TradeIT/models"
	"container/heap"
	"errors"
	"fmt"
	"math"
	"sync"
)

type Orderbook struct {
	orderTable  map[uint64]*Node
	buy_orders  map[float64]*DoublyLinkedList
	sell_orders map[float64]*DoublyLinkedList
	bids_prices MaxHeap
	asks_prices MinHeap
	tradeCount  int64
	buyCount    int64
	sellCount   int64
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
		orderTable:  make(map[uint64]*Node),
	}
	heap.Init(&ob.bids_prices)
	heap.Init(&ob.asks_prices)
	return &ob
}

func (ob *Orderbook) Lock()   { ob.mu.Lock() }   //used while testing
func (ob *Orderbook) Unlock() { ob.mu.Unlock() } //used while testing

func (ob *Orderbook) InsertOrder(orderData models.Metadata) error {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.internalInsertOrder(orderData)
}

func (ob *Orderbook) internalInsertOrder(orderData models.Metadata) error {
	if orderData.Side == "buy" {
		if _, exists := ob.buy_orders[orderData.Price]; !exists {
			ob.buy_orders[orderData.Price] = &DoublyLinkedList{}
			heap.Push(&ob.bids_prices, orderData.Price) // inserts the price in Heap
		}
		ob.orderTable[orderData.Id] = ob.buy_orders[orderData.Price].PushBack(orderData) // insert order details in buy order list and also stores the nodes reference in the orderTable
		ob.buyCount++
		return nil
	} else if orderData.Side == "sell" {
		if _, exists := ob.sell_orders[orderData.Price]; !exists {
			ob.sell_orders[orderData.Price] = &DoublyLinkedList{}
			heap.Push(&ob.asks_prices, orderData.Price) // inserts the price in Heap
		}
		ob.orderTable[orderData.Id] = ob.sell_orders[orderData.Price].PushBack(orderData) //insert order details in buy order list
		ob.sellCount++
		return nil
	} else {
		return errors.New("invalid order type")
	}
}

func (ob *Orderbook) Matcher(order models.Metadata) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	if order.Side == "buy" { //buy order
		for ob.asks_prices.Len() > 0 && order.Remq > 0 {
			bestAsk := ob.asks_prices.Peek()
			if order.Price < bestAsk { //buyers price is
				break
			}

			askList, exists := ob.sell_orders[bestAsk]
			if !exists {
				heap.Pop(&ob.asks_prices)
			} else {
				for node := askList.Head; node != nil && order.Remq > 0; node = node.Next {
					matchQuantity := min(node.Metadata.Remq, order.Remq)
					order.Remq -= matchQuantity
					node.Metadata.Remq -= matchQuantity
					executionPrice := math.Round(((order.Price+node.Metadata.Price)/2.0)*100) / 100
					node.Metadata.AvgPrice = float64(matchQuantity) * executionPrice
					order.AvgPrice = float64(matchQuantity) * executionPrice
					ob.tradeCount++
					if node.Metadata.Remq == 0 {
						if node.Metadata.Remq == 0 {
							// remove the node's reference from the orderTable
							delete(ob.orderTable, node.Metadata.Id)
							// remove the node from the linked list
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
							node.Metadata.AvgPrice = node.Metadata.AvgPrice / float64(node.Metadata.Quantity)
							askList.Size--
							ob.sellCount--
							//send the message to the user that order has been filled
						}
					}
					if order.Remq == 0 {
						order.AvgPrice = order.AvgPrice / float64(order.Quantity) //finalises the avg price
						delete(ob.orderTable, order.Id)                           //removes the order from the orderTable
						//send the message to the user that the order has been filled
					}

				}
				if askList.Size == 0 {
					delete(ob.sell_orders, bestAsk)
					heap.Pop(&ob.asks_prices)
				}
			}
		}
		if order.Remq > 0 { //adds the order to the orderbook if not filled
			ob.internalInsertOrder(order)
		}
	} else {
		// Proccessing the sell order

		for ob.bids_prices.Len() > 0 && order.Remq > 0 { //checking for the bids available
			bestBid := ob.bids_prices.Peek()
			if order.Price > bestBid { //breaks if bids are less than the asks
				break
			}
			bidList, exists := ob.buy_orders[bestBid]
			if !exists {
				heap.Pop(&ob.bids_prices)
			} else {
				for node := bidList.Head; node != nil && order.Remq > 0; node = node.Next {
					//trade occurs
					matchQuantity := min(order.Remq, node.Metadata.Remq)
					order.Remq -= matchQuantity
					node.Metadata.Remq -= matchQuantity
					executionPrice := math.Round(((order.Price+node.Metadata.Price)/2.0)*100) / 100
					node.Metadata.AvgPrice = float64(matchQuantity) * executionPrice
					order.AvgPrice = float64(matchQuantity) * executionPrice
					ob.tradeCount++
					//register the trade here
					//
					if node.Metadata.Remq == 0 { //buyer order is fullfilled and node has become stale
						if node.Metadata.Remq == 0 {
							// remove the node's reference from the orderTable
							delete(ob.orderTable, node.Metadata.Id)
							// remove the node from the linked list
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
							node.Metadata.AvgPrice = node.Metadata.AvgPrice / float64(node.Metadata.Quantity)
							bidList.Size--
							ob.buyCount--
							//send the message to the user that the order has been filled
						}
					}

					if order.Remq == 0 {
						order.AvgPrice = order.AvgPrice / (float64(order.Quantity))
						delete(ob.orderTable, order.Id) //removes the order from the orderTable
						//send the message to the user that the order has been filled
					}
				}
				if bidList.Size == 0 {
					if _, exists := ob.buy_orders[bestBid]; exists {
						delete(ob.buy_orders, bestBid)
					}
					heap.Pop(&ob.bids_prices)
				}
			}
		}
		if order.Remq > 0 {
			ob.internalInsertOrder(order)
		}

	}

	if ob.asks_prices.Len() == 0 && ob.bids_prices.Len() == 0 {
		//close the matcher and orderbook
		fmt.Println("Time to close the orderbook and matcher")
	}

}

func (ob *Orderbook) DisplayResult() {
	fmt.Printf("Trades: %d |\nBuyOrders: %d|\nSellOrders: %d|\nTotal: %d", ob.tradeCount, ob.buyCount, ob.sellCount, (ob.tradeCount + ob.sellCount + ob.buyCount))
	fmt.Println("\nHeap lengths: ")
	fmt.Printf("bids: %d\nasks: %d", ob.bids_prices.Len(), ob.asks_prices.Len())
	fmt.Println("")
}

func (ob *Orderbook) CancelOrder(orderId uint64) error {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	node, exists := ob.orderTable[orderId]
	if !exists {
		return errors.New("order has been proceessed or does not exist")
	}
	difference := node.Metadata.Quantity - node.Metadata.Remq
	if node.Metadata.Side == "buy" {
		sig := ob.buy_orders[node.Metadata.Price].RemoveNode(node)
		if sig != nil {
			delete(ob.buy_orders, node.Metadata.Price) // deletes the entry of linkedlist in the map
		}
		if difference != 0 {
			node.Metadata.Remq = 0
			node.Metadata.Quantity = difference
			//send the user the message of the partial filled quantity order is success
			ob.buyCount--
			return errors.New("order was partially filled,rest of order has been cancelled")
		}
		ob.buyCount--
	} else if node.Metadata.Side == "sell" {
		sig := ob.sell_orders[node.Metadata.Price].RemoveNode(node)
		if sig != nil {
			delete(ob.sell_orders, node.Metadata.Price)
		}
		if difference != 0 {
			node.Metadata.Remq = 0
			node.Metadata.Quantity = difference
			//send the user the message of the partial filled quantity order is success
			ob.sellCount--
			return errors.New("order was partially filled,rest of order has been cancelled")
		}
		ob.sellCount--
	} else {
		return errors.New("unsupported order type")
	}
	delete(ob.orderTable, orderId) //deletes the node and order entry from the ordertable
	return nil
}

func (ob *Orderbook) ModifyQuantity(orderId uint64, newQuantity int) error {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	order_node, exists := ob.orderTable[orderId]
	if !exists {
		return errors.New("order has been processed in the mean time or order doesn't exist")
	}
	if newQuantity <= 0 {
		return errors.New("invalid quantity (negative/zero quantity)")
	}
	difference := newQuantity - (order_node.Metadata.Quantity - order_node.Metadata.Remq)

	if difference < 0 { //Partially filled more than than newQuantity
		return errors.New("order has been partially filled and quantity can't be decreased ")
	} else if difference == 0 { //order has just completed and rest of the order is not to be filled
		if order_node.Metadata.Side == "buy" {

		} else if order_node.Metadata.Side == "sell" {
			sig := ob.sell_orders[order_node.Metadata.Price].RemoveNode(order_node)
			if sig != nil {
				delete(ob.sell_orders, order_node.Metadata.Price)
			}
			ob.sellCount--
		}
		order_node.Metadata.Quantity = newQuantity
	} else {
		order_node.Metadata.Quantity = newQuantity
		order_node.Metadata.Remq = difference
	}
	return nil
}

func (ob *Orderbook) ModifyPrice(orderId uint64, newPrice float64) error {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	order_node, exists := ob.orderTable[orderId]
	if !exists {
		return errors.New("order has been processed or order doesn't exists")
	}
	if newPrice <= 0 {
		return errors.New("invalid price (negative or zero price)")
	}
	if newPrice == order_node.Metadata.Price {
		return nil
	}
	if order_node.Metadata.Side == "buy" {
		sig := ob.buy_orders[order_node.Metadata.Price].RemoveNode(order_node)
		if sig != nil {
			delete(ob.buy_orders, order_node.Metadata.Price) // deletes the entry of linkedlist in the map
		}
		ob.buyCount--

	} else if order_node.Metadata.Side == "sell" {
		sig := ob.sell_orders[order_node.Metadata.Price].RemoveNode(order_node)
		if sig != nil {
			delete(ob.sell_orders, order_node.Metadata.Price)
		}
		ob.sellCount--
	} else {
		return errors.New("invalid order side")
	}
	delete(ob.orderTable, orderId)
	order_node.Metadata.Price = newPrice
	ob.Unlock()
	ob.Matcher(order_node.Metadata)
	ob.Lock()
	return nil
}
