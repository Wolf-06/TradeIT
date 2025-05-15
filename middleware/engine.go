package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func redis_test() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_addr"),
		Password: os.Getenv("redis_password"),
		DB:       0,
	})

	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ping)

	err = client.Set(context.Background(), "name", "TradeIT", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Get(context.Background(), "name").Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}

func main() {
	redis_test()
}
