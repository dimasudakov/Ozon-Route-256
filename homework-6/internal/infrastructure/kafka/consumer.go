package kafka

import (
	"github.com/IBM/sarama"
	"time"
)

type Consumer struct {
	brokers        []string
	SingleConsumer sarama.Consumer
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, config)

	if err != nil {
		return nil, err
	}

	return &Consumer{
		brokers:        brokers,
		SingleConsumer: consumer,
	}, nil
}
