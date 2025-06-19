package engine

import (
	"TradeIT/database"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type OrderEngine struct {
	redisClient *redis.Client
	db          *gorm.DB
}

func InitOrderEngine() *OrderEngine {
	return &OrderEngine{
		redisClient: GetEngineClient(),
		db:          database.SetDB(),
	}
}

func (engine *OrderEngine) ProcessOrderQueue() {
	for {
		queueData, err := engine.redisClient.BLPop(context.Background(), 0, "order").Result()
		if err != nil {
			fmt.Println("Error in processing the order queue: ", err)
			continue
		}
		err = engine.ProcessOrder(queueData)
		if err != nil {
			fmt.Println("Error in processing the order: ", err)
			continue
		}
	}
}

func (engine *OrderEngine) ProcessOrder(orderdata any) error {
	fmt.Println(orderdata)
	return nil
}

func EngineTest() {
	engine := InitOrderEngine()
	engine.ProcessOrderQueue()
}
