package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/user"
	"github.com/SQUASHD/hbit/user/userdb"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	connectionStr := os.Getenv("USER_DB_URL")
	db, err := hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to user database: %s", err)
	}

	err = hbit.DBMigrateUp(db, hbit.MigrationData{
		FS:      user.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})
	if err != nil {
		log.Fatalf("failed to run migration of user database: %v", err)
	}

	queries := userdb.New(db)
	userSvc := user.NewService(db, queries)

	userRouter := http.NewUserRouter(userSvc)
	wrappedRouter := http.ChainMiddleware(
		userRouter,
	)
	server, err := http.NewServer(
		wrappedRouter,
		http.WithServerOptionsPortFromEnv("USER_SVC_PORT"),
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
