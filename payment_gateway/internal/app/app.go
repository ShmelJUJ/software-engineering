package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/config"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/subscriber"
	"github.com/ShmelJUJ/software-engineering/pkg/kafka"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
)

const (
	defaultMessageFetchBytes  = 1024 * 1024
	defaultAutoCommitEnabled  = true
	defaultAutoCommitInterval = time.Second
)

func getSubscriberSaramaConfig() (*sarama.Config, error) {
	saramaConfig := sarama.NewConfig()

	saramaVersion, err := sarama.ParseKafkaVersion("3.0.1")
	if err != nil {
		return nil, err
	}

	saramaConfig.Version = saramaVersion
	saramaConfig.Consumer.Fetch.Default = defaultMessageFetchBytes
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = defaultAutoCommitEnabled
	saramaConfig.Consumer.Offsets.AutoCommit.Interval = defaultAutoCommitInterval

	return saramaConfig, nil
}

// Run is a function which run application.
func Run(cfg *config.Config) {
	ctx := context.Background()

	l, err := logger.NewLogrusLogger(cfg.LoggerCfg.Level)
	if err != nil {
		log.Fatal("failed to create logger: ", err)
	}

	saramaCfg, err := getSubscriberSaramaConfig()
	if err != nil {
		l.Fatal("failed to get subscriber sarama config", map[string]interface{}{
			"error": err,
		})
	}

	kafkaSubscriber, err := kafka.NewSubscriber(
		cfg.KafkaSubscriberCfg.Brokers,
		kafka.WithSubscriberSaramaConfig(saramaCfg),
		kafka.WithSubscriberInitTopic(&sarama.TopicDetail{
			NumPartitions:     int32(cfg.KafkaSubscriberCfg.TopicDetails.NumPartitions),
			ReplicationFactor: int16(cfg.KafkaSubscriberCfg.TopicDetails.ReplicationFactor),
		}),
		kafka.WithSubscriberConsumerGroup("payment_gateway"),
	)
	if err != nil {
		l.Fatal("failed to create kafka subscriber", map[string]interface{}{
			"error": err,
		})
	}

	kafkaRouter, err := kafka.NewBrokerRouter()
	if err != nil {
		l.Fatal("failed to create kafka router", map[string]interface{}{
			"error": err,
		})
	}

	kafkaPublisher, err := kafka.NewPublisher(cfg.KafkaPublisherCfg.Brokers)
	if err != nil {
		l.Fatal("failed to create kafka publisher", map[string]interface{}{
			"error": err,
		})
	}

	sub, err := subscriber.NewTransactionSubscriber(
		cfg.KafkaSubscriberCfg.SubscriberCfg,
		l,
		kafkaRouter,
		kafkaSubscriber,
		kafkaPublisher,
		cfg.KafkaPublisherCfg.PublisherCfg,
		cfg.AlgorandCfg,
	)
	if err != nil {
		l.Fatal("failed to create new transaction subscriber", map[string]interface{}{
			"error": err,
		})
	}

	sub.RegisterCancelledTransactionHandler()
	sub.RegisterProcessedTransactionHandler()

	go func() {
		if err := sub.Run(ctx); err != nil {
			l.Fatal("failed to run subscriber", map[string]interface{}{
				"error": err,
			})
		}
	}()

	defer func() {
		if err := sub.Stop(); err != nil {
			l.Error("failed to stop subscriber", map[string]interface{}{
				"error": err,
			})
		}
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	l.Info("Graceful shutdown...", map[string]interface{}{
		"signal": s.String(),
	})
}
