package engine

import "TradeIT/models"

// -----------| HEAPS |------------------

// MinHeap implementation
type MinHeap []float64

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(float64))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Peek returns the minimum element without removing it
func (h *MinHeap) Peek() float64 {
	if h.Len() == 0 {
		return 0
	}
	return (*h)[0]
}

// MaxHeap implementation using embedding approach
type MaxHeap struct {
	MinHeap
}

func (h MaxHeap) Less(i, j int) bool {
	return h.MinHeap[i] > h.MinHeap[j]
}

// Peek returns the maximum element without removing it
func (h *MaxHeap) Peek() float64 {
	if h.Len() == 0 {
		return 0
	}
	return h.MinHeap[0]
}

//---------DOUBLY-LINKED-LIST-------------

type Node struct {
	Metadata models.Metadata
	Next     *Node
	Prev     *Node
}

type DoublyLinkedList struct {
	Head *Node
	Tail *Node
	Size int
}

func (dll *DoublyLinkedList) PushFront(order models.Metadata) {
	newNode := &Node{Metadata: order}
	if dll.Head == nil {
		dll.Head = newNode
		dll.Tail = newNode
	} else {
		newNode.Next = dll.Head
		dll.Head.Prev = newNode
		dll.Head = newNode
	}
	dll.Size++
}

func (dll *DoublyLinkedList) PushBack(order models.Metadata) {
	newNode := &Node{Metadata: order}
	if dll.Head == nil {
		dll.Head = newNode
		dll.Tail = newNode
	} else {
		newNode.Prev = dll.Tail
		dll.Tail.Next = newNode
		dll.Tail = newNode
	}
	dll.Size++
}

func (dll *DoublyLinkedList) RemoveFront() {
	dll.Head = dll.Head.Next
	dll.Head.Prev = nil
	dll.Size--
	if dll.Size == 0 {
		dll.Tail = nil
	}
}
