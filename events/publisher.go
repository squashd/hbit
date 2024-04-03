package events

import (
	"github.com/SQUASHD/hbit/config"
	"github.com/wagslane/go-rabbitmq"
)

func NewPublisher(rabbitmqConf config.RabbitMQ) (*rabbitmq.Publisher, *rabbitmq.Conn, error) {
	connStr := config.NewRabbitConnectionString(rabbitmqConf)
	conn, err := rabbitmq.NewConn(connStr)
	if err != nil {
		return nil, nil, err
	}
	authPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, nil, err
	}
	return authPub, conn, nil
}
