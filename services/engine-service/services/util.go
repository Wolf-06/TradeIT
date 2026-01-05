package services

import (
	"math/rand"
	"time"
)

func createUserId() uint64 {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return uint64(random.Intn(9000) + 1000)
}
