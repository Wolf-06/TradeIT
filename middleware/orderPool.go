package middleware

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
			return new(models.Metadata)
		},
	}
	return &op
}

func (op *OrderPool) AcquireOrder() *models.Metadata {
	return op.pool.Get().(*models.Metadata)
}

func (op *OrderPool) ReleaseOrder(o *models.Metadata) {
	*o = models.Metadata{}
	op.pool.Put(o)
}
