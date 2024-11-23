package di

import (
	"log/slog"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type kafkaConsumer struct {
	consumer sarama.ConsumerGroup
}

type kafkaConsumerOptions struct {
	username          string
	password          string
	sessionTimeout    time.Duration
	heartbeatInterval time.Duration
	bufferSize        int
	maxRetry          int
	brokers           []string
	group             string
}

const (
	defaultKafkaSessionTimeout    = 10 * time.Second
	defaultKafkaHeartbeatInterval = 3 * time.Second

	defaultKafkaConsumerBufferSize = 256
	defaultKafkaConsumerMaxRetry   = 3
)

func newKafkaConsumer(opts kafkaConsumerOptions) (*kafkaConsumer, error) {
	sessionTimeout := opts.sessionTimeout
	if sessionTimeout == 0 {
		sessionTimeout = defaultKafkaSessionTimeout
	}

	heartbeatInterval := opts.heartbeatInterval
	if heartbeatInterval == 0 {
		heartbeatInterval = defaultKafkaHeartbeatInterval
	}

	bufferSize := opts.bufferSize
	if bufferSize == 0 {
		bufferSize = defaultKafkaConsumerBufferSize
	}

	maxRetry := opts.maxRetry
	if maxRetry == 0 {
		maxRetry = defaultKafkaConsumerMaxRetry
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Session.Timeout = sessionTimeout
	saramaConfig.Consumer.Group.Heartbeat.Interval = heartbeatInterval
	saramaConfig.ChannelBufferSize = bufferSize
	saramaConfig.Consumer.Group.Rebalance.Retry.Max = maxRetry
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = true
	saramaConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	saramaConfig.Net.SASL.User = opts.username
	saramaConfig.Net.SASL.Password = opts.password

	consumer, err := sarama.NewConsumerGroup(opts.brokers, opts.group, saramaConfig)
	if err != nil {
		return nil, err
	}

	slog.Info("[di.newKafkaConsumer] Kafka consumer started",
		slog.String("brokers", strings.Join(opts.brokers, ",")),
	)

	return &kafkaConsumer{consumer}, nil
}
