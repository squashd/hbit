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
	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/task/taskdb"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	connectionStr := os.Getenv("TASK_DB_URL")
	db, err := hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to task database: %s", err)
	}

	err = hbit.DBMigrateUp(db, hbit.MigrationData{
		FS:      task.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})
	if err != nil {
		log.Fatalf("failed to run migration of task database: %v", err)
	}

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	publisher, conn, err := events.NewPublisher(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create task publisher: %s", err)
	}
	defer conn.Close()

	queries := taskdb.New(db)
	taskSvc := task.NewService(db, queries, publisher)

	consumer, conn, err := events.NewTaskEventConsumer(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create task consumer: %s", err)
	}
	defer conn.Close()

	taskRouter := http.NewTaskRouter(taskSvc)

	wrappedRouter := http.ChainMiddleware(
		taskRouter,
	)
	server, err := http.NewServer(
		wrappedRouter,
		http.WithServerOptionsPortFromEnv("TASK_SVC_PORT"),
	)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
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

		if err := taskSvc.CleanUp(); err != nil {
			log.Fatalf("Task service cleanup failure: %v", err)
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
			events.TaskMessageHandler(taskSvc),
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
