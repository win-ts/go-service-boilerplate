package di

import (
	"log/slog"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

type kafkaProducerOptions struct {
	username string
	password string
	brokers  []string
	timeout  time.Duration
	maxRetry int
}

const (
	defaultKafkaProducerTimeout  = 3 * time.Second
	defaultKafkaProducerMaxRetry = 3
)

func newKafkaProducer(opts kafkaProducerOptions) (*kafkaProducer, error) {
	timeout := opts.timeout
	if timeout == 0 {
		timeout = defaultKafkaProducerTimeout
	}

	maxRetry := opts.maxRetry
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
	saramaConfig.Net.SASL.User = opts.username
	saramaConfig.Net.SASL.Password = opts.password

	p, err := sarama.NewSyncProducer(opts.brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	slog.Info("[di.newKafkaProducer] Kafka producer started",
		slog.String("brokers", strings.Join(opts.brokers, ",")),
	)

	return &kafkaProducer{
		producer: p,
	}, nil
}
