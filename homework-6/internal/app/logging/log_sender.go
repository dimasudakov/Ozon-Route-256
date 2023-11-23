//go:generate mockgen -source=./logger.go -destination=../mocks/logger.go -package=app_mock

package logging

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/infrastructure/kafka"
	"time"
)

type Logger interface {
	Log(msg LogMessage)
	Warning(msg LogMessage)
	Error(msg LogMessage)
}

type LogMessage struct {
	Time        time.Time
	LogLevel    string
	RequestURI  string
	RequestType string
	Method      string
	Body        string
	Info        string
}

type KafkaLogger struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaLogger(producer *kafka.Producer, topic string) *KafkaLogger {
	return &KafkaLogger{
		producer: producer,
		topic:    topic,
	}
}

func (k *KafkaLogger) Log(msg LogMessage) {
	msg.LogLevel = "INFO"
	k.sendLog(msg)
}

func (k *KafkaLogger) Warning(msg LogMessage) {
	msg.LogLevel = "WARNING"
	k.sendLog(msg)
}

func (k *KafkaLogger) Error(msg LogMessage) {
	msg.LogLevel = "ERROR"
	k.sendLog(msg)
}

func (k *KafkaLogger) sendLog(msg LogMessage) error {
	msg.Time = time.Now()
	producerMsg, err := k.buildMessage(msg)
	if err != nil {
		return err
	}

	_, _, err = k.producer.SendSyncMessage(producerMsg)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaLogger) buildMessage(msg LogMessage) (*sarama.ProducerMessage, error) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     k.topic,
		Value:     sarama.ByteEncoder(jsonMsg),
		Timestamp: time.Now(),
	}, nil
}
