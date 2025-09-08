package engine

import (
	"TradeIT/database"
	"TradeIT/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type OrderEngine struct {
	ledger          *Ledger
	ord_redisClient *redis.Client
	mod_redisClient *redis.Client
	db              *gorm.DB
}

var cancellation = make(map[uint64]struct{})
var modification = make(map[uint64]struct{})

func InitOrderEngine() *OrderEngine {
	InitTrade()
	TestDB()
	return &OrderEngine{
		ledger:          InitLedger(),
		ord_redisClient: GetEngineClient(0),
		mod_redisClient: GetEngineClient(1),
		db:              database.SetDB(),
	}
}

func (engine *OrderEngine) ProcessOrderQueue() {
	var order models.Order
	fmt.Printf("Starting the OrderQueue\n")
	for {
		queueData, err := engine.ord_redisClient.BLPop(context.Background(), 0, "IncomingOrder").Result()
		if err != nil {
			fmt.Println("Error in processing the order queue: ", err)
			continue
		}
		err = json.Unmarshal([]byte(queueData[1]), &order)
		if err != nil {
			fmt.Println("Error in umarshalling the order from the redis server: ", err)
		}
		if _, exists := cancellation[order.Id]; exists {
			continue
		}
		fmt.Println("processing order of id: ", order.Id)
		err = engine.ledger.ProcessOrder(order)
		if err != nil {
			fmt.Println("Error has occured in the matcher: ", err)
			//resolve the error
		}
		// order = models.Order{} // resets the data, to avoid redundancy
	}
}

func (engine *OrderEngine) ProcessCancellationQueue() {
	for {
		//recieves the orders that are too be
		queueData, err := engine.mod_redisClient.BLPop(context.Background(), 0, "cancel").Result()
		if err != nil {
			log.Println("Error in processing the cancellation queue: ", err)
			continue
		}
		//data is tranformed to the required types
		data := strings.Split(string(queueData[1]), ":")
		orderid, _ := strconv.ParseUint(data[0], 10, 64)

		result := engine.ledger.CancelOrder(data[1], orderid)
		if err = engine.mod_redisClient.LPush(context.Background(), "cancel:"+data[0], result).Err(); err != nil {
			log.Println("Error in publishing the status of cancelling order: ")
		}
	}
}

func EngineTest() {
	engine := InitOrderEngine()
	engine.ProcessOrderQueue()
}
