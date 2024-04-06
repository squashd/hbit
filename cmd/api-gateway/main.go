package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/auth/authdb"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/events"
	"github.com/SQUASHD/hbit/http"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	jwtConf := config.NewJwtConfig(
		config.WithJwtOptionsSecretFromEnv("JWT_SECRET"),
	)

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	publisher, conn, err := events.NewPublisher(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create auth publisher: %s", err)
	}
	defer conn.Close()

	db, err := auth.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to auth database: %v", err)
	}

	queries := authdb.New(db)
	authSvc := auth.NewService(publisher, jwtConf, db, queries)

	gateway, err := http.SetUpAPIGateway(authSvc, jwtConf)
	if err != nil {
		log.Fatalf("cannot create gateway: %s", err)
	}

	server, err := http.NewServer(
		http.ChainMiddleware(
			gateway,
			http.CORSMiddleware,
			http.LoggerMiddleware,
		),
		http.WithServerOptionsPortFromEnv("API_GATEWAY_PORT"),
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

		ctx, cancel := context.WithTimeout(context.Background(), server.IdleTimeout)
		defer cancel()
		if err := authSvc.Cleanup(); err != nil {
			log.Fatalf("Auth service cleanup failure: %v", err)
		}

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
