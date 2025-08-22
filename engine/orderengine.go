package engine

import (
	"TradeIT/database"
	"TradeIT/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type OrderEngine struct {
	ledger      Ledger
	redisClient *redis.Client
	db          *gorm.DB
}

func InitOrderEngine() *OrderEngine {
	return &OrderEngine{
		ledger:      *InitLedger(),
		redisClient: GetEngineClient(),
		db:          database.SetDB(),
	}
}

func (engine *OrderEngine) ProcessOrderQueue() {
	var order models.Metadata
	for {
		queueData, err := engine.redisClient.BLPop(context.Background(), 0, "order").Result()
		if err != nil {
			fmt.Println("Error in processing the order queue: ", err)
			continue
		}
		err = json.Unmarshal([]byte(queueData[1]), &order)
		if err != nil {
			fmt.Println("Error in umarshalling the order from the redis server: ", err)
		}
		go engine.ProcessOrder(order)
		order = models.Metadata{} // resets the data, to avoid redundancy
	}
}

func (engine *OrderEngine) ProcessOrder(orderdata models.Metadata) {
	err := engine.ledger.ProcessOrder(orderdata)
	if err != nil {
		fmt.Println("Error has occured in the matcher: ", err)
		//resolve the error
	}
}

func EngineTest() {
	engine := InitOrderEngine()
	engine.ProcessOrderQueue()
}
