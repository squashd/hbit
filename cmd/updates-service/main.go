package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/events"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	consumer, conn, err := events.NewUpdateEventConsumer(config.RabbitMQ{})
	if err != nil {
		log.Fatalf("cannot create feat consumer: %s", err)
	}
	defer conn.Close()
	eventHandler := events.NewUpdateConsumerHandler()

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		consumer.Close()

		close(closed)
	}()
	fmt.Println("Starting update consumer")
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := consumer.Run(
			eventHandler.HandleEvents,
		); err != nil {
			log.Fatalf("cannot start consuming: %s", err)
		}
	}()
	wg.Wait()

	<-closed
	log.Println("Server closed")
}
