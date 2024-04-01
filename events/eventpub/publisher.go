package eventpub

import (
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/events/messagebroker"
	"github.com/wagslane/go-rabbitmq"
)

type (
	Publisher interface {
		Publish(message []byte, routingKey string, headers map[string]any) error
	}

	PublisherOptions struct {
		ExchangeName string
		ExchangeType string
	}

	rabbitmqpub struct {
		publisher *rabbitmq.Publisher
	}
)

func NewRabbitMQPublisher(opts PublisherOptions, conf config.RabbitmqConnnection) (Publisher, func(), error) {
	connStr := messagebroker.NewRabbitConnectionString(conf)
	conn, err := rabbitmq.NewConn(connStr)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		conn.Close()
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsExchangeName(opts.ExchangeName),
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, cleanup, err
	}

	return &rabbitmqpub{publisher: publisher}, cleanup, nil
}

func (p *rabbitmqpub) Publish(message []byte, routingKey string, payload map[string]any) error {
	return p.publisher.Publish(
		message,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsHeaders(payload),
	)
}
