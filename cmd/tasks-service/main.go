package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit/events"
	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/task"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	publisher, conn, err := events.NewTaskPublisher()
	if err != nil {
		log.Fatalf("cannot create task publisher: %s", err)
	}
	defer conn.Close()

	taskDb, err := task.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to task database: %v", err)
	}
	taskRepo := task.NewRepository(taskDb)
	taskSvc := task.NewService(taskRepo, publisher)

	taskRouter := http.NewTaskRouter(taskSvc)
	server, err := http.NewServer(
		taskRouter,
		http.WithServerOptionsPort(9001),
	)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		fmt.Println("\nShutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), server.IdleTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failure: %v", err)
		}

		close(closed)
	}()
	fmt.Printf("Server is running on port %s\n", server.Addr)
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start server: %s", err)
	}
}
