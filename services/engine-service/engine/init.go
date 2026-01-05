package engine

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var EngineClient *redis.Client

func InitEngineClient(db int) {
	EngineClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_addr"),
		Password: os.Getenv("redis_password"),
		DB:       db,
	})
}

func GetEngineClient(db int) *redis.Client {
	if EngineClient == nil {
		InitEngineClient(db)
		fmt.Println("successfully created the redis client")
	}
	return EngineClient
}

func custom(a, b int) {
	var mt = make(map[string]string)
	fmt.Println(mt)
}
