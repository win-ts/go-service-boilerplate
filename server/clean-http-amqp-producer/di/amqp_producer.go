package di

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type amqpProducer struct {
	channel *amqp.Channel
}

type amqpProducerOptions struct {
	host            string
	username        string
	password        string
	queueName       string
	queueDurable    bool
	queueAutoDelete bool
	queueExclusive  bool
	queueNoWait     bool
}

func newAMQPProducer(opts amqpProducerOptions) (*amqpProducer, error) {
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

	slog.Info("[di.newAMQPProducer] rabbitmq producer started",
		slog.Any("queue_name", opts.queueName),
	)

	return &amqpProducer{ch}, nil
}
