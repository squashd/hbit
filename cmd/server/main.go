package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SQUASHD/hbit/achievement"
	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/character"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/http"
	"github.com/SQUASHD/hbit/quest"
	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/user"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/wagslane/go-rabbitmq"
)

func main() {

	serverConfig, jwtConfig, err := setupConfig()
	if err != nil {
		log.Fatalf("cannot setup config: %s", err)
	}

	authSvc, authCleanup, err := setUpAuthService(jwtConfig)
	if err != nil {
		log.Fatalf("cannot create auth service: %s", err)
	}
	defer authCleanup()

	charSvc, charCleanup, err := setUpCharacterService()
	if err != nil {
		log.Fatalf("cannot create character service: %s", err)
	}
	defer charCleanup()

	questSvc, questCleanup, err := setUpQuestService()
	if err != nil {
		log.Fatalf("cannot create quest service: %s", err)
	}
	defer questCleanup()

	achievementSvc, taskSubscriber, achCleanup, err := setUpAchievementService(config.RabbitmqConfig{})
	if err != nil {
		log.Fatalf("cannot create achievement service: %s", err)
	}
	defer achCleanup()

	taskSvc, taskCleanup, err := setUpTaskService()
	if err != nil {
		log.Fatalf("cannot create task service: %s", err)
	}
	defer taskCleanup()

	userSvc, userCleanup, err := setUpUserService()
	if err != nil {
		log.Fatalf("cannot create user service: %s", err)
	}
	defer userCleanup()

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

		close(closed)
	}()

	go func() {
		if err := taskSubscriber.Run(
			func(d rabbitmq.Delivery) rabbitmq.Action {
				msg, err := json.Marshal(d.Body)
				if err != nil {
					log.Printf("cannot marshal message: %s", err)
					return rabbitmq.NackDiscard
				}
				switch d.Type {
				case "task_done":
					log.Printf("task done event: %s", msg)
				default:
					log.Printf("unknown event: %s", d.Type)
					return rabbitmq.NackDiscard
				}
				return rabbitmq.Ack
			},
		); err != nil {
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

func setUpAuthService(jwtConf config.JwtConfig) (auth.Service, func(), error) {
	rabbitmqConf := config.RabbitmqConfig{}
	connStr := config.NewRabbitConnectionString(rabbitmqConf)
	conn, err := rabbitmq.NewConn(connStr)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		conn.Close()
	}
	authPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("task_updates"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, cleanup, err
	}
	authDb, err := auth.NewDatabase()
	if err != nil {
		return nil, cleanup, err
	}

	cleanup = func() {
		authDb.Close()
		cleanup()
	}

	authRepo := auth.NewRepository(authDb)
	authSvc := auth.NewService(authRepo, jwtConf, authPub)
	return authSvc, cleanup, nil
}

func setUpCharacterService() (character.Service, func(), error) {
	charDb, err := character.NewDatabase()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		charDb.Close()
	}
	charRepo := character.NewRepository(charDb)
	charSvc := character.NewService(charRepo)
	return charSvc, cleanup, nil
}

func setUpQuestService() (quest.Service, func(), error) {
	questDb, err := quest.NewDatabase()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		questDb.Close()
	}
	questRepo := quest.NewRepository(questDb)
	questSvc := quest.NewService(questRepo)
	return questSvc, cleanup, nil
}

func setUpAchievementService(rabbitmqConf config.RabbitmqConfig) (achievement.Service, *rabbitmq.Consumer, func(), error) {
	connStr := config.NewRabbitConnectionString(rabbitmqConf)
	conn, err := rabbitmq.NewConn(connStr)

	cleanup := func() {
		conn.Close()
	}
	if err != nil {
		return nil, nil, cleanup, err
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"my_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	achievementDb, err := achievement.NewDatabase()
	if err != nil {
		return nil, nil, cleanup, err
	}
	achievementRepo := achievement.NewRepository(achievementDb)
	achievementSvc := achievement.NewService(achievementRepo)
	return achievementSvc, consumer, cleanup, nil
}

func setUpTaskService() (task.Service, func(), error) {
	rabbitmqConf := config.RabbitmqConfig{}
	connStr := config.NewRabbitConnectionString(rabbitmqConf)
	conn, err := rabbitmq.NewConn(connStr)
	cleanup := func() {
		conn.Close()
	}
	if err != nil {
		return nil, cleanup, err
	}
	taskPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("task_updates"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)

	taskDb, err := task.NewDatabase()
	if err != nil {
		return nil, cleanup, err
	}

	cleanup = func() {
		taskDb.Close()
		cleanup()
	}

	taskRepo := task.NewRepository(taskDb)
	taskSvc := task.NewService(taskRepo, taskPub)
	return taskSvc, cleanup, nil
}

func setUpUserService() (user.Service, func(), error) {
	userDb, err := user.NewDatabase()
	if err != nil {
		log.Fatalf("cannot create user database: %s", err)
	}
	cleanup := func() {
		userDb.Close()
	}
	userRepo := user.NewReposiory(userDb)
	userSvc := user.NewService(userRepo)
	return userSvc, cleanup, nil
}

func setupConfig() (config.ServerConfig, config.JwtConfig, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return config.ServerConfig{}, config.JwtConfig{}, fmt.Errorf("JWT_SECRET is required")
	}
	jwtConfig := config.NewJwtConfig(jwtSecret, "access", "refresh", 1*60, 24*60*60)

	serverConfig, err := config.NewServerConfigFromEnv()
	if err != nil {
		return config.ServerConfig{}, config.JwtConfig{}, err
	}

	return serverConfig, jwtConfig, nil
}
