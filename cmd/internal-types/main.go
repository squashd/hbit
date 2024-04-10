package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/events"
	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/task/taskdb"
)

func main() {
	connectionStr := os.Getenv("RPG_DB_URL")
	db, err := hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: connectionStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to rpg database: %s", err)
	}
	bool := os.Getenv("DEBUG")
	fmt.Println("DEBUG: ", bool)

	err = hbit.DBMigrateUp(db, hbit.MigrationData{
		FS:      rpg.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})
	if err != nil {
		log.Fatalf("failed to run migration of rpg database: %v", err)
	}

	queries := rpgdb.New(db)

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	publisher, conn, err := events.NewPublisher(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create rpg publisher: %s", err)
	}
	defer conn.Close()
	charPublisher, conn2, err := events.NewPublisher(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create rpg publisher: %s", err)
	}
	defer conn2.Close()
	taskConnStr := os.Getenv("TASK_DB_URL")
	taskDb, err := hbit.NewDatabase(hbit.NewDbParams{
		ConnectionStr: taskConnStr,
		Driver:        hbit.DbDriverLibsql,
	})
	if err != nil {
		log.Fatalf("cannot connect to rpg database: %s", err)
	}
	err = hbit.DBMigrateUp(taskDb, hbit.MigrationData{
		FS:      task.Migrations,
		Dialect: "sqlite",
		Dir:     "schemas",
	})
	if err != nil {
		log.Fatalf("failed to run migration of rpg database: %v", err)
	}

	taskQueries := taskdb.New(taskDb)

	taskPublisher, conn3, err := events.NewPublisher(rabbitmqUrl)
	if err != nil {
		log.Fatalf("cannot create rpg publisher: %s", err)
	}
	defer conn3.Close()
	taskSvc := task.NewService(taskDb, taskQueries, taskPublisher)

	questSvc := quest.NewService(db, queries)
	charSvc := character.NewService(db, queries, charPublisher)
	rpgSvc := rpg.NewService(rpg.NewServiceParams{
		QuestSvc:     questSvc,
		CharacterSvc: charSvc,
		Queries:      queries,
		Publisher:    publisher,
		Db:           db,
	})
	typesRouter := http.NewTypesRouter(rpgSvc, questSvc, charSvc, taskSvc)
	server, err := http.NewServer(
		typesRouter,
		http.WithServerOptionsPort(9500),
	)
	if err != nil {
		panic(err)
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
	fmt.Println("Server is running on port 9500")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

}
