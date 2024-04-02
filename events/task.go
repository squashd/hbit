package events

import (
	"github.com/SQUASHD/hbit/config"
	"github.com/wagslane/go-rabbitmq"
)

func NewTaskPublisher() (*rabbitmq.Publisher, *rabbitmq.Conn, error) {
	rabbitmqConf := config.RabbitMQ{}
	connStr := config.NewRabbitConnectionString(rabbitmqConf)
	conn, err := rabbitmq.NewConn(connStr)
	if err != nil {
		return nil, nil, err
	}
	taskPublisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, nil, err
	}

	return taskPublisher, conn, nil
}
