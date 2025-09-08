package repo

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
			return new(models.Order)
		},
	}
	return &op
}

func (op *OrderPool) AcquireOrder() *models.Order {
	return op.pool.Get().(*models.Order)
}

func (op *OrderPool) ReleaseOrder(o *models.Order) {
	*o = models.Order{}
	op.pool.Put(o)
}
