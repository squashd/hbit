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
	"github.com/SQUASHD/hbit/feat"
	"github.com/SQUASHD/hbit/feat/featdb"
	"github.com/SQUASHD/hbit/http"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	db, err := feat.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to feat database: %v", err)
	}

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	publisher, conn, err := events.NewPublisher(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create publisher: %s", err)
	}
	defer conn.Close()

	queries := featdb.New(db)
	featSvc := feat.NewService(db, queries, publisher)

	consumer, conn, err := events.NewFeatEventConsumer(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create feat consumer: %s", err)
	}
	defer conn.Close()
	eventHandler := events.NewFeatEventHandler(featSvc)

	featRouter := http.NewFeatRouter(featSvc)
	wrappedRouter := http.ChainMiddleware(
		featRouter,
		http.CORSMiddleware,
		http.LoggerMiddleware,
	)
	server, err := http.NewServer(
		wrappedRouter,
		http.WithServerOptionsPort(80),
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
