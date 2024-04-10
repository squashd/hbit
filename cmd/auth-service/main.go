package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/auth/authdb"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/events"
	"github.com/SQUASHD/hbit/http"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {

	connectionStr := os.Getenv("TASK_DB_URL")
	db, err := hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to auth database: %s", err)
	}

	err = hbit.DBMigrateUp(db, hbit.MigrationData{
		FS:      auth.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})
	if err != nil {
		log.Fatalf("failed to run migration of auth database: %v", err)
	}

	jwtConf := config.NewJwtConfig(
		config.WithJwtOptionsSecretFromEnv("JWT_SECRET"),
	)

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	fmt.Printf("rabbitmqUrl: %s\n", rabbitmqUrl)
	publisher, conn, err := events.NewPublisher(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create auth publisher: %s", err)
	}
	defer conn.Close()

	queries := authdb.New(db)

	authSvc := auth.NewService(publisher, jwtConf, db, queries)

	authRouter := http.NewAuthRouter(authSvc, jwtConf)
	wrappedRouter := http.ChainMiddleware(
		authRouter,
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
