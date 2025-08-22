package engine

import (
	"TradeIT/models"
	"sync"
)

type TradePool struct {
	pool sync.Pool
}

func InitTradePool() *TradePool {
	var tp TradePool
	tp.pool = sync.Pool{
		New: func() interface{} {
			return new(models.TradeDetails)
		},
	}
	return &tp
}

func (tp *TradePool) acquireTrade() *models.TradeDetails {
	return tp.pool.Get().(*models.TradeDetails)
}

func (tp *TradePool) releaseTrade(td *models.TradeDetails) {
	*td = models.TradeDetails{}
	tp.pool.Put(td)
}
