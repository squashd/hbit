package events

import (
	"github.com/SQUASHD/hbit/config"
	"github.com/wagslane/go-rabbitmq"
)

func NewAuthPublisher() (*rabbitmq.Publisher, *rabbitmq.Conn, error) {
	rabbitmqConf := config.RabbitMQ{}
	connStr := config.NewRabbitConnectionString(rabbitmqConf)
	conn, err := rabbitmq.NewConn(connStr)
	if err != nil {
		return nil, nil, err
	}
	authPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, nil, err
	}
	return authPub, conn, nil
}
