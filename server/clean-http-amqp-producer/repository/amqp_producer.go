package repository

import (
	"context"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/pkg/utils"
)

type amqpProducerRepository struct {
	channel   *amqp.Channel
	queueName string
}

// AMQPProducerRepositoryConfig represents the configuration of AMQP Producer repository
type AMQPProducerRepositoryConfig struct {
	QueueName string
}

// AMQPProducerRepositoryDependencies represents the dependencies of AMQP Producer repository
type AMQPProducerRepositoryDependencies struct {
	Channel *amqp.Channel
}

// NewAMQPProducerRepository creates a new instance of AMQPProducerRepository
func NewAMQPProducerRepository(c AMQPProducerRepositoryConfig, d AMQPProducerRepositoryDependencies) AMQPProducerRepository {
	return &amqpProducerRepository{
		channel:   d.Channel,
		queueName: c.QueueName,
	}
}

// ProduceJSONMessage produces a JSON message to the AMQP
func (r *amqpProducerRepository) ProduceJSONMessage(ctx context.Context, event dto.Event) error {
	if err := r.channel.PublishWithContext(ctx,
		"",          // exchange
		r.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(utils.EncodeJSONtoBytes(event)),
		}); err != nil {
		return errors.Wrap(err, "[amqpProducerRepository.ProduceJSONMessage] error publishing message")
	}

	return nil
}
