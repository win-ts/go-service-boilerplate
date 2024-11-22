// Package kafka provides the Kafka producer
package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

// Producer represents the Kafka producer
type Producer struct {
	Producer sarama.SyncProducer
}

// ProducerOptions represents the options for the Kafka producer
type ProducerOptions struct {
	Username string
	Password string
	Brokers  []string
	Timeout  time.Duration
	MaxRetry int
}

const (
	defaultKafkaProducerTimeout  = 3 * time.Second
	defaultKafkaProducerMaxRetry = 3
)

// NewProducer creates a new Kafka producer
func NewProducer(opts ProducerOptions) (*Producer, error) {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = defaultKafkaProducerTimeout
	}

	maxRetry := opts.MaxRetry
	if maxRetry == 0 {
		maxRetry = defaultKafkaProducerMaxRetry
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Timeout = timeout
	saramaConfig.Producer.Retry.Max = maxRetry
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	saramaConfig.Net.SASL.User = opts.Username
	saramaConfig.Net.SASL.Password = opts.Password

	p, err := sarama.NewSyncProducer(opts.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: p,
	}, nil
}

// Close closes the Kafka producer
func (k *Producer) Close() error {
	return k.Producer.Close()
}
