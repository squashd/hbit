package events

import (
	"github.com/wagslane/go-rabbitmq"
)

func NewPublisher(url string) (*rabbitmq.Publisher, *rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(url)
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
