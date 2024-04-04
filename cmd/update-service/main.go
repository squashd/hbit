package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SQUASHD/hbit/events"
	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/updates"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	consumer, conn, err := events.NewUpdateEventConsumer(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create feat consumer: %s", err)
	}
	defer conn.Close()
	eventHandler := events.NewUpdateConsumerHandler()

	svc := updates.NewService()

	router := http.NewUpdatesRouter(svc)

	server, err := http.NewServer(
		http.ChainMiddleware(
			router,
			http.CORSMiddleware,
		),
		http.WithServerOptionsPortFromEnv("UPDATES_SVC_PORT"),
	)
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
	}
	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		consumer.Close()

		fmt.Println("\nShutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), server.IdleTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failure: %v", err)
		}

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
	fmt.Printf("Server is running on port %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start server: %s", err)
	}
	wg.Wait()

	<-closed
	log.Println("Server closed")
}
