package eventsub

import (
	"context"

	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/events/messagebroker"
	"github.com/wagslane/go-rabbitmq"
)

type (
	Subscriber interface {
		StartConsuming(ctx context.Context, handlerFunc rabbitmq.Handler) error
	}

	SubscriberOptions struct {
		ExchangeName string
		ExchangeType string
		QueueName    string
		RoutingKey   string
	}

	rabbitmqSubscriber struct {
		consumer *rabbitmq.Consumer
	}
)

func NewSubscriber(opts SubscriberOptions, conf config.RabbitmqConnnection) (Subscriber, func(), error) {
	connStr := messagebroker.NewRabbitConnectionString(conf)
	conn, err := rabbitmq.NewConn(connStr)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		conn.Close()
	}

	consumer, err := rabbitmq.NewConsumer(
		conn,
		opts.QueueName,
		rabbitmq.WithConsumerOptionsRoutingKey(opts.RoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(opts.ExchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, cleanup, err
	}

	return &rabbitmqSubscriber{consumer: consumer}, cleanup, nil
}

func (r *rabbitmqSubscriber) StartConsuming(ctx context.Context, handlerFunc rabbitmq.Handler) error {
	return r.consumer.Run(
		handlerFunc,
	)
}
