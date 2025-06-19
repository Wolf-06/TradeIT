package engine

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var EngineClient *redis.Client

func InitEngineClient() {
	EngineClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_addr"),
		Password: os.Getenv("redis_password"),
		DB:       0,
	})
}

func GetEngineClient() *redis.Client {
	if EngineClient == nil {
		InitEngineClient()
	}
	return EngineClient
}
