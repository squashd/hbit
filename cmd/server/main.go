package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit/achievement"
	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/character"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/events/eventhandler"
	"github.com/SQUASHD/hbit/events/eventpub"
	"github.com/SQUASHD/hbit/events/eventsub"
	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/quest"
	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/user"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {

	// Set up JWT Config
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	jwtConfig := config.NewJwtConfig(jwtSecret, "access", "refresh", 1*60, 24*60*60)

	// Set up Server Config
	serverConfig, err := config.NewServerConfigFromEnv()
	if err != nil {
		log.Fatalf("cannot create server config: %s", err)
	}

	// Set up Auth Service
	authDb, err := auth.NewDatabase()
	if err != nil {
		log.Fatalf("cannot create auth database: %s", err)
	}
	authRepo := auth.NewRepository(authDb)
	authSvc := auth.NewService(authRepo, jwtConfig)

	// Set up Character Service
	charDb, err := character.NewDatabase()
	if err != nil {
		log.Fatalf("cannot create character database: %s", err)
	}
	charRepo := character.NewRepository(charDb)
	charSvc := character.NewService(charRepo)

	// Set up Quest Service
	questDb, err := quest.NewDatabase()
	if err != nil {
		log.Fatalf("cannot create quest database: %s", err)
	}
	questRepo := quest.NewRepository(questDb)
	questSvc := quest.NewService(questRepo)

	// Set up Achievement Service
	achievementDb, err := achievement.NewDatabase()
	if err != nil {
		log.Fatalf("cannot create achievement database: %s", err)
	}
	achievementRepo := achievement.NewRepository(achievementDb)
	achievementSvc := achievement.NewService(achievementRepo)

	// Set up Task Service
	rabbitmqConf := config.RabbitmqConnnection{}
	pubOpts := eventpub.PublisherOptions{
		ExchangeName: "task_updates",
		ExchangeType: "topic",
	}
	taskEventPublisher, taskPublisherCleanup, err := eventpub.NewRabbitMQPublisher(pubOpts, rabbitmqConf)
	if err != nil {
		log.Fatalf("cannot create rabbitmq publisher: %s", err)
	}

	taskDb, err := task.NewDatabase()
	if err != nil {
		log.Fatalf("cannot create task database: %s", err)
	}
	taskRepo := task.NewRepository(taskDb)
	taskSvc := task.NewService(taskRepo, taskEventPublisher)

	taskHandler := eventhandler.TaskHandler(taskSvc)
	subOpts := eventsub.SubscriberOptions{
		ExchangeName: "task_updates",
		ExchangeType: "topic",
		QueueName:    "",
		RoutingKey:   "task.*",
	}
	taskSubscriber, tasksubscriberCleanup, err := eventsub.NewSubscriber(subOpts, rabbitmqConf)
	if err != nil {
		log.Fatalf("cannot create rabbitmq subscriber: %s", err)
	}

	// Set up User Service
	userDb, err := user.NewDatabase()
	if err != nil {
		log.Fatalf("cannot create user database: %s", err)
	}
	userRepo := user.NewReposiory(userDb)
	userSvc := user.NewService(userRepo)

	monolith := http.NewServerMonolith(serverConfig, jwtConfig, authSvc, charSvc, questSvc, achievementSvc, taskSvc, userSvc)
	server, err := http.NewServer(serverConfig, monolith.RegisterRoutes())
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		fmt.Println("\nShutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), serverConfig.TimeoutIdle)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failure: %v", err)
		}

		taskPublisherCleanup()
		tasksubscriberCleanup()

		if err = authDb.Close(); err != nil {
			log.Fatalf("Database shutdown failure: %v", err)
		}
		if err = taskDb.Close(); err != nil {
			log.Fatalf("Database shutdown failure: %v", err)
		}
		if err = questDb.Close(); err != nil {
			log.Fatalf("Database shutdown failure: %v", err)
		}
		if err = charDb.Close(); err != nil {
			log.Fatalf("Database shutdown failure: %v", err)
		}
		if err = achievementDb.Close(); err != nil {
			log.Fatalf("Database shutdown failure: %v", err)
		}
		if err = userDb.Close(); err != nil {
			log.Fatalf("Database shutdown failure: %v", err)
		}

		close(closed)
	}()

	go func() {
		if err := taskSubscriber.StartConsuming(context.Background(), taskHandler); err != nil {
			log.Fatalf("cannot start consuming: %s", err)
		}
	}()

	fmt.Printf("Server is running on port %d\n", serverConfig.Port)
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start server: %s", err)
	}

	<-closed
	log.Println("Server closed")

}
