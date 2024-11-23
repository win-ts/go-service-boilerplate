package repository

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-kafka-producer/dto"
)

type kafkaProducerRepository struct {
	producer sarama.SyncProducer
	topic    string
}

// KafkaProducerRepositoryConfig represents the configuration for the producer repository
type KafkaProducerRepositoryConfig struct {
	TopicName string
}

// KafkaProducerRepositoryDependencies represents the dependencies for the producer repository
type KafkaProducerRepositoryDependencies struct {
	Producer sarama.SyncProducer
}

// NewKafkaProducerRepository creates a new producer repository
func NewKafkaProducerRepository(c KafkaProducerRepositoryConfig, d KafkaProducerRepositoryDependencies) KafkaProducerRepository {
	return &kafkaProducerRepository{
		producer: d.Producer,
		topic:    c.TopicName,
	}
}

// Produce produces a message to the Kafka topic
func (r *kafkaProducerRepository) Produce(message dto.Event) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	partition, offset, err := r.producer.SendMessage(&sarama.ProducerMessage{
		Topic: r.topic,
		Key:   sarama.StringEncoder("exampleKey"),
		Value: sarama.StringEncoder(msg),
	})
	if err != nil {
		return err
	}

	slog.Info("[producerRepository.Produce] message produced",
		slog.Int("partition", int(partition)),
		slog.Int64("offset", offset),
	)
	return nil
}
