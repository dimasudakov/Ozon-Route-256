package kafka

import (
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type Producer struct {
	brokers      []string
	syncProducer sarama.SyncProducer
}

func newSyncLogger(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()

	syncProducerConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll
	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true
	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "error occured during creating sync infrastructure producer")
	}

	return syncProducer, nil
}

func NewProducer(brokers []string) (*Producer, error) {
	syncProducer, err := newSyncLogger(brokers)
	if err != nil {
		return nil, err
	}

	return &Producer{
		brokers:      brokers,
		syncProducer: syncProducer,
	}, nil
}

func (k *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.syncProducer.SendMessage(message)
}

func (k *Producer) Close() error {
	err := k.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "infrastructure.Connector.Close")
	}
	return nil
}
