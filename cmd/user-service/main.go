package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/user"
	"github.com/SQUASHD/hbit/user/userdb"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	db, err := user.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to user database: %v", err)
	}

	quries := userdb.New(db)
	userSvc := user.NewService(db, quries)

	userRouter := http.NewUserRouter(userSvc)
	wrappedRouter := http.ChainMiddleware(
		userRouter,
		http.CORSMiddleware,
		http.LoggerMiddleware,
	)
	server, err := http.NewServer(
		wrappedRouter,
		http.WithServerOptionsPort(80),
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
