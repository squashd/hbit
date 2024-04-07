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
	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	db, err := rpg.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to rpg database: %v", err)
	}

	queries := rpgdb.New(db)

	questSvc := quest.NewService(db, queries)
	characterSvc := character.NewService(db, queries)

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	publisher, conn, err := events.NewPublisher(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create rpg publisher: %s", err)
	}
	defer conn.Close()

	rpgSvc := rpg.NewService(publisher, queries, db)

	consumer, conn, err := events.NewRPGEventConsumer(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create rpg consumer: %s", err)
	}
	defer conn.Close()

	rpgRouter := http.NewRPGRouter(characterSvc, questSvc, rpgSvc)
	wrappedRouter := http.ChainMiddleware(
		rpgRouter,
	)
	server, err := http.NewServer(
		wrappedRouter,
		http.WithServerOptionsPortFromEnv("RPG_SVC_PORT"),
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
			events.RPGMessageHandler(rpgSvc),
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
