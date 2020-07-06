package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

// KafkaProducerHelper represents kafka producer helper interface
type KafkaProducerHelper interface {
	Send(ctx context.Context, producerName, topic, key string, value interface{}) error
}

// asyncKafkaProducer represents async kafka producer
type asyncKafkaProducer struct {
	producerName     string
	producerInstance sarama.AsyncProducer
	ready            bool
}

// syncKafkaProducer represents sync kafka producer
type syncKafkaProducer struct {
	producerName     string
	producerInstance sarama.SyncProducer
	ready            bool
}

// NewKafkaProducerHelper creates an instance
func NewKafkaProducerHelper(isAsync bool, producerName string, brokers []string, version string) KafkaProducerHelper {
	if isAsync {
		asyncKafkaProducer, err := initAsyncKafkaProducer(producerName, brokers, version)
		if err != nil {
			zap.S().Panic("Failed to init async Kafka producer", zap.Error(err))
		}
		return asyncKafkaProducer
	}
	syncKafkaProducer, err := initSyncKafkaProducer(producerName, brokers, version)
	if err != nil {
		zap.S().Panic("Failed to init sync Kafka producer", zap.Error(err))
	}
	return syncKafkaProducer
}

// initSyncKafkaProducer represents init sync kafka producer
func initSyncKafkaProducer(producerName string, brokers []string, version string) (*syncKafkaProducer, error) {
	kafkaVersion, _ := sarama.ParseKafkaVersion(version)
	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		zap.S().Fatalf("Can't init kafka producer with err %v", zap.Error(err))
	}
	zap.S().Infof("Init sync Kafka Producer successfully")
	syncKafkaProducer := &syncKafkaProducer{
		producerName:     producerName,
		producerInstance: producer,
		ready:            true,
	}
	return syncKafkaProducer, nil
}

// initAsyncKafkaProducer represents init async kafka producer
func initAsyncKafkaProducer(producerName string, brokers []string, version string) (*asyncKafkaProducer, error) {
	kafkaVersion, _ := sarama.ParseKafkaVersion(version)
	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		zap.S().Fatalf("Can't init kafka producer with err %v", zap.Error(err))
	}
	zap.S().Infof("Init async Kafka Producer successfully")
	go func() {
		for {
			select {
			case err := <-producer.Errors():
				{
					zap.S().Errorw(fmt.Sprintf("Failed while sending message %s to topic %s", err.Msg.Key, err.Msg.Topic), zap.Error(err.Err))
					producer.Input() <- err.Msg
				}
			case sucess := <-producer.Successes():
				{
					zap.S().Infow(fmt.Sprintf("%s - Sent message successfully to topic %s", sucess.Key, sucess.Topic))
				}
			}
		}
	}()
	asyncKafkaProducer := &asyncKafkaProducer{
		producerName:     producerName,
		producerInstance: producer,
		ready:            true,
	}
	return asyncKafkaProducer, nil
}

// Send represents AsyncKafkaProducer send
func (h *asyncKafkaProducer) Send(ctx context.Context, producerName string, topic string, key string, value interface{}) error {

	if !h.ready {
		return errors.New("Producer is not ready at all")
	}

	if h.producerName != producerName {
		return errors.New("Wrong producer name")
	}

	buffer, err := json.Marshal(value)
	if err != nil {
		return errors.New("Can't marshal object")
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(buffer),
	}
	zap.S().Debug("Send to queue")
	h.producerInstance.Input() <- message
	return nil
}

// Send represents SyncKafkaProducer Send
func (h *syncKafkaProducer) Send(ctx context.Context, producerName string, topic string, key string, value interface{}) error {
	return nil
}
