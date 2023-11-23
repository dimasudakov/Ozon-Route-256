package logging

import (
	"context"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/infrastructure/kafka"
	"sync"
)

type HandleFunc func(message *sarama.ConsumerMessage, ch chan *sarama.ConsumerMessage)

type LogReceiver struct {
	consumer       *kafka.Consumer
	messageHandler HandleFunc
}

func NewLogReceiver(consumer *kafka.Consumer, messageHandler HandleFunc) *LogReceiver {
	return &LogReceiver{
		consumer:       consumer,
		messageHandler: messageHandler,
	}
}

func (r *LogReceiver) Subscribe(ctx context.Context, topic string, ch chan *sarama.ConsumerMessage) (*sync.WaitGroup, error) {
	var wg sync.WaitGroup

	partitionList, err := r.consumer.SingleConsumer.Partitions(topic)
	if err != nil {
		return nil, err
	}

	initialOffset := sarama.OffsetNewest

	wg.Add(len(partitionList))
	for _, partition := range partitionList {
		pc, err := r.consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			return nil, err
		}
		go func(pc sarama.PartitionConsumer, partition int32) {
			defer func() {
				pc.Close()
				wg.Done()
			}()
			for {
				select {
				case <-ctx.Done():
					return
				case message, ok := <-pc.Messages():
					if !ok {
						return
					}
					r.messageHandler(message, ch)
				}
			}
		}(pc, partition)
	}

	return &wg, nil
}
