package infrustructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/logging"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/infrastructure/kafka"
	"log"
	"sync"
	"testing"
	"time"
)

var (
	testTopicName = "logs_test"
	msgHandler    = func(msg *sarama.ConsumerMessage, msgChan chan *sarama.ConsumerMessage) {
		msgChan <- msg
	}
)

type KafkaLogsTestTopic struct {
	Logger   logging.Logger
	Receiver *logging.LogReceiver
	msgChan  chan *sarama.ConsumerMessage
	wg       *sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	sync.Mutex
}

func NewKafkaLogsTestTopic(logger logging.Logger, consumer kafka.Consumer) *KafkaLogsTestTopic {
	return &KafkaLogsTestTopic{
		Logger:   logger,
		Receiver: logging.NewLogReceiver(&consumer, msgHandler),
		msgChan:  make(chan *sarama.ConsumerMessage),
	}
}

func (k *KafkaLogsTestTopic) SetUp(t *testing.T, chanSize int) {
	t.Helper()
	k.Lock()
	k.msgChan = make(chan *sarama.ConsumerMessage, chanSize)
	k.ctx, k.cancel = context.WithCancel(context.Background())
	var err error
	k.wg, err = k.Receiver.Subscribe(k.ctx, testTopicName, k.msgChan)
	if err != nil {
		panic(err)
	}
}

func (k *KafkaLogsTestTopic) TearDown() {
	defer k.Unlock()
	// ждем пока у всех сообщений в кафке истечет retention.ms и они удалятся сами
	time.Sleep(1 * time.Second)
}

func (k *KafkaLogsTestTopic) CreateTestTopic(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Error creating cluster admin: %v", err)
		return err
	}
	defer admin.Close()

	topic := "logs_test"
	retentionMs := "200"
	detail := sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 2,
		ConfigEntries: map[string]*string{
			"retention.ms": &retentionMs,
		},
	}

	err = admin.CreateTopic(topic, &detail, false)
	if err != nil {
		if !errors.Is(err, sarama.ErrTopicAlreadyExists) {
			log.Fatalf("Error creating topic: %v", err)
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (k *KafkaLogsTestTopic) CheckLogs(t *testing.T, messages []logging.LogMessage) bool {
	// ждем пока все логи запишутся в кафку
	time.Sleep(1 * time.Second)

	k.cancel()
	k.wg.Wait()
	close(k.msgChan)
	checkResult := true
	idx := 0
	for message := range k.msgChan {
		if idx >= len(messages) {
			idx++
			break
		}
		logMsg := logging.LogMessage{}
		err := json.Unmarshal(message.Value, &logMsg)
		if err != nil {
			fmt.Println("Error during json unmarshalling")
		}
		if diff := cmp.Diff(messages[idx], logMsg, cmpopts.IgnoreFields(logging.LogMessage{},
			"Time", "Method", "RequestURI", "Info")); diff != "" {
			t.Errorf("Differences between logs: (-expected +actual):\n%s", diff)
			checkResult = false
		}
		idx++
	}

	if idx != len(messages) {
		t.Errorf("Expected logs number: %d, received: %d", len(messages), idx)
		checkResult = false
	}
	return checkResult
}
