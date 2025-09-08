package services

import (
	"TradeIT/models"
	"sync"
)

type OrderPool struct {
	pool sync.Pool
}

func InitOrderPool() *OrderPool {
	var op OrderPool
	op.pool = sync.Pool{
		New: func() interface{} {
			return new(models.MetaOrder)
		},
	}
	return &op
}

func (op *OrderPool) acquireOrder() *models.MetaOrder {
	return op.pool.Get().(*models.MetaOrder)
}

func (op *OrderPool) releaseOrder(o *models.MetaOrder) {
	*o = models.MetaOrder{}
	op.pool.Put(o)
}
