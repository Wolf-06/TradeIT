package services

import (
	"math/rand"
	"time"
)

func createUserId() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Intn(9000) + 1000
}
