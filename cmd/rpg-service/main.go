package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/events"
	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	rpgDb, err := rpg.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to rpg database: %v", err)
	}
	characterRepo := character.NewRepository(rpgDb)
	questRepo := quest.NewRepository(rpgDb)
	questSvc := quest.NewService(questRepo)
	characterSvc := character.NewService(characterRepo)

	publisher, conn, err := events.NewPublisher()
	if err != nil {
		log.Fatalf("cannot create rpg publisher: %s", err)
	}
	defer conn.Close()

	rpgSvc := rpg.NewService(characterSvc, questSvc, publisher)

	consumer, conn, err := events.NewRPGEventConsumer(config.RabbitMQ{})
	if err != nil {
		log.Fatalf("cannot create rpg consumer: %s", err)
	}
	defer conn.Close()
	eventHandler := events.NewRPGConsumerHandler(rpgSvc)

	rpgRouter := http.NewRPGRouter(characterSvc, questSvc)
	server, err := http.NewServer(
		rpgRouter,
		http.WithServerOptionsPort(8080),
	)
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		fmt.Println("\nShutting down server...")

		publisher.Close()
		consumer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), server.IdleTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failure: %v", err)
		}

		close(closed)
	}()
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
	fmt.Printf("Server is running on port %s\n", server.Addr)
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start server: %s", err)
	}
	wg.Wait()

	<-closed
	log.Println("Server closed")
}
