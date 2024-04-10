package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/events"
	"github.com/SQUASHD/hbit/feat"
	"github.com/SQUASHD/hbit/feat/featdb"
	"github.com/SQUASHD/hbit/http"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	connectionStr := os.Getenv("ACH_DB_URL")
	db, err := hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to feat database: %s", err)
	}

	err = hbit.DBMigrateUp(db, hbit.MigrationData{
		FS:      feat.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})
	if err != nil {
		log.Fatalf("failed to run migration of feat database: %v", err)
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
	featRouter := http.NewFeatRouter(featSvc)
	wrappedRouter := http.ChainMiddleware(
		featRouter,
	)
	server, err := http.NewServer(
		wrappedRouter,
		http.WithServerOptionsPortFromEnv("FEAT_SVC_PORT"),
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

		if err := featSvc.CleanUp(); err != nil {
			log.Fatalf("Feat service cleanup failure: %v", err)
		}

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
			events.FeatsMessageHandler(featSvc),
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
