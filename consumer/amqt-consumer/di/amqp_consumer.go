package di

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type amqpConsumer struct {
	deliveryChan <-chan amqp.Delivery
}

type amqpConsumerOptions struct {
	host              string
	username          string
	password          string
	queueName         string
	queueDurable      bool
	queueAutoDelete   bool
	queueExclusive    bool
	queueNoWait       bool
	consumerName      string
	consumerAutoAck   bool
	consumerExclusive bool
	consumerNoLocal   bool
	consumerNoWait    bool
}

func newAMQPConsumer(opts amqpConsumerOptions) (*amqpConsumer, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", opts.username, opts.password, opts.host))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		opts.queueName,
		opts.queueDurable,
		opts.queueAutoDelete,
		opts.queueExclusive,
		opts.queueNoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	deliveryChan, err := ch.Consume(
		opts.queueName,
		opts.consumerName,
		opts.consumerAutoAck,
		opts.consumerExclusive,
		opts.consumerNoLocal,
		opts.consumerNoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	slog.Info("[di.newAMQPConsumer] rabbitmq consumer connected",
		slog.String("queue_name", opts.queueName),
	)

	return &amqpConsumer{deliveryChan}, nil
}
