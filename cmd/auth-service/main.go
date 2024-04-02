package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/events"
	"github.com/SQUASHD/hbit/http"
)

func main() {
	jwtConf := config.NewJwtConfig(
		config.WithJwtOptionsSecretFromEnv("JWT_SECRET"),
	)
	publisher, conn, err := events.NewAuthPublisher()
	if err != nil {
		log.Fatalf("cannot create auth publisher: %s", err)
	}
	defer conn.Close()

	authDb, err := auth.NewDatabase()
	authRepo := auth.NewRepository(authDb)
	authSvc := auth.NewService(authRepo, jwtConf, publisher)

	authRouter := http.NewAuthRouter(authSvc)
	server, err := http.NewServer(authRouter)
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
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
	fmt.Printf("Server is running on port %d\n", server.Addr)
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start server: %s", err)
	}
}
