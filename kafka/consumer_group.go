package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

// KafkaConsumerGroup represents KafkaConsumerGroup
type KafkaConsumerGroup struct {
	ready     chan bool
	groups    sarama.ConsumerGroup
	MessageCh chan *sarama.ConsumerMessage
	ErrorCh   chan *sarama.ConsumerError
}

// InitConsumerGroup represents initConsumerGroup
func InitConsumerGroup(ctx context.Context, topics []string, group string, brokers []string, version string, rebalanceStrategy string, isOldest bool, clientID string) (*KafkaConsumerGroup, error) {
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		zap.S().Panicw("Error while parsing kafka version", zap.Error(err))
	}
	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Consumer.Return.Errors = true
	// config.ClientID = clientID
	switch rebalanceStrategy {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		zap.S().Panicw(fmt.Sprintf("Can't initialize rebalance strategy of consumer group: %v", rebalanceStrategy))
	}
	if isOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	kafkaConsumer := &KafkaConsumerGroup{
		ready:     make(chan bool),
		MessageCh: make(chan *sarama.ConsumerMessage),
		ErrorCh:   make(chan *sarama.ConsumerError),
	}

	consumer, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		zap.S().Panicw(fmt.Sprintf("Error while creating consumer group: %v", zap.Error(err)))
	}
	kafkaConsumer.groups = consumer

	go func() {
		for {
			if err := consumer.Consume(ctx, topics, kafkaConsumer); err != nil {
				if err == sarama.ErrClosedConsumerGroup {
					zap.S().Info("Consumer Group is closed")
					break
				}
				zap.S().Errorw(fmt.Sprintf("Consumer Group: Failed to consume: %v", zap.Error(err)))
				time.Sleep(time.Second)
			}
			if ctx.Err() != nil {
				return
			}
			kafkaConsumer.ready = make(chan bool)
		}
	}()
	<-kafkaConsumer.ready
	return kafkaConsumer, nil
}

// HandleCloseConsumeGroup represents HandleCloseConsumeGroup
func (kfg *KafkaConsumerGroup) HandleCloseConsumeGroup() {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	go func(sigterm chan os.Signal) {
		<-sigterm

		zap.S().Info("*******STOP KAFKA CONSUMER **********")
		if err := kfg.groups.Close(); err != nil {
			zap.S().Errorw(fmt.Sprintf("Error closing client: %v", err))
		}
	}(sigterm)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (kfp *KafkaConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	close(kfp.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (kfp *KafkaConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim represents ConsumeClaim
func (kfg *KafkaConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		zap.S().Infow(fmt.Sprintf("%s - Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Key), string(message.Value), message.Timestamp, message.Topic))
		session.MarkMessage(message, "")
		kfg.MessageCh <- message
	}
	return nil
}
