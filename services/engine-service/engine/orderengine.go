package engine

import (
	"TradeIT/internal/database"
	"TradeIT/internal/models"
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
<<<<<<< HEAD:services/engine-service/engine/orderengine.go
		//recieves the orders that are too be cancelled.
		fmt.Println("Waiting for the cancellation request.")
		queueData, err := engine.mod_redisClient.BLPop(context.Background(), 0, "CancelQueue").Result()
		fmt.Println("Recieved Request")
=======
		//recieves the orders that are too be
		queueData, err := engine.mod_redisClient.BLPop(context.Background(), 0, "cancel").Result()
>>>>>>> parent of 6daa9c4 (Completed the Order Creation and Matching Pipeline.):engine/orderengine.go
		if err != nil {
			log.Println("Error in processing the cancellation queue: ", err)
			continue
		}
<<<<<<< HEAD:services/engine-service/engine/orderengine.go
		fmt.Println(queueData)
=======
		//data is tranformed to the required types
>>>>>>> parent of 6daa9c4 (Completed the Order Creation and Matching Pipeline.):engine/orderengine.go
		data := strings.Split(string(queueData[1]), ":")
		orderid, _ := strconv.ParseUint(data[0], 10, 64)

		result := engine.ledger.CancelOrder(data[1], orderid)

		if result != nil {
			switch result.Error() {
			case "ORDER_NOT_IN_MATCHER":
				if !checkTrade(orderid) {
					cancellation[orderid] = struct{}{}
					fmt.Println("Cancelled Successfully")
				} else {
				}

			case "ORDER_PARTIAL_CANCELLED":

			case "ORDER_UNSUPPORTED_TYPE":
				log.Println("Cancel failed: Unsupported order type.")

			default:
				log.Println("Cancel failed: ", result.Error())
			}
		} else {
			log.Println("Cancel success: Order cancelled.")
		}
	}
}

<<<<<<< HEAD:services/engine-service/engine/orderengine.go
func (engine *OrderEngine) ProccessQuantityModification() {
	fmt.Println("---Starting the Quantity Modification Request Queue---")
	for {
		queueData, err := engine.mod_redisClient.BLPop(context.Background(), 0, "PriceModifyQueue").Result()
		fmt.Println("Recieved Request")
		if err != nil {
			log.Println("Error in processing the cancellation queue: ", err)
			continue
		}
		fmt.Println(queueData)
		data := strings.Split(string(queueData[1]), ":")
		orderid, _ := strconv.ParseUint(data[0], 10, 64)

		_ = engine.ledger.CancelOrder(data[1], orderid)
	}
=======
func EngineTest() {
	engine := InitOrderEngine()
	engine.ProcessOrderQueue()
>>>>>>> parent of 6daa9c4 (Completed the Order Creation and Matching Pipeline.):engine/orderengine.go
}
