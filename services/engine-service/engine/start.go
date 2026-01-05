package engine

import (
	"fmt"
	"log"
)

func start() {
	engine := InitOrderEngine()
	fmt.Println("Starting the engine")
	log.Println("Started the cancellation processing")
	go engine.ProcessCancellationQueue()
	log.Println("starting the order Processing ")
	engine.ProcessOrderQueue()
}
