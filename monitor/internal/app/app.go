package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/ShmelJUJ/software-engineering/monitor/internal/broker/publisher"
	"github.com/ShmelJUJ/software-engineering/monitor/internal/broker/subscriber"
	"github.com/ShmelJUJ/software-engineering/monitor/internal/generated/restapi"
	"github.com/ShmelJUJ/software-engineering/monitor/internal/generated/restapi/operations"

	apiMonitor "github.com/ShmelJUJ/software-engineering/monitor/internal/generated/restapi/operations/monitor"

	"github.com/go-openapi/loads"

	"github.com/ShmelJUJ/software-engineering/monitor/config"
	"github.com/ShmelJUJ/software-engineering/monitor/internal/api/handler"
	"github.com/ShmelJUJ/software-engineering/pkg/kafka"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	gen "github.com/ShmelJUJ/software-engineering/user/gen"
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

func Run(cfg *config.Config) {
	ctx := context.Background()

	l, err := logger.NewLogrusLogger("info")
	if err != nil {
		log.Fatalln("failed to create new logger: ", err)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		l.Fatal("failed to get swagger spec", map[string]interface{}{
			"error": err,
		})
	}

	userClient, err := gen.NewClient(cfg.UserClientCfg.URL)
	if err != nil {
		l.Fatal("failed to make new user client", map[string]interface{}{
			"error": err,
		})
	}

	monitorHandler := handler.NewMonitorHandler(l, userClient)

	api := operations.NewMonitorAPI(swaggerSpec)

	api.MonitorProcessHandler = apiMonitor.ProcessHandlerFunc(monitorHandler.ProcessHandler)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			l.Error("failed to shutdown server", map[string]interface{}{
				"error": err,
			})
		}
	}()

	server.ConfigureAPI()

	go func() {
		server.Port = cfg.HTTPCfg.Port
		if err := server.Serve(); err != nil {
			l.Fatal("failed to serve server", map[string]interface{}{
				"error": err,
			})
		}
	}()

	kafkaPublisher, err := kafka.NewPublisher(cfg.PublisherCfg.Brokers)
	if err != nil {
		l.Fatal("failed to create kafka publisher", map[string]interface{}{
			"error": err,
		})
	}

	monitorPublisher := publisher.NewMonitorPublisher(
		l,
		kafkaPublisher,
	)

	saramaCfg, err := getSubscriberSaramaConfig()
	if err != nil {
		l.Fatal("failed to get subscriber sarama config", map[string]interface{}{
			"error": err,
		})
	}

	// Run subscriber
	kafkaSubscriber, err := kafka.NewSubscriber(
		cfg.SubscriberCfg.Brokers,
		kafka.WithSubscriberSaramaConfig(saramaCfg),
		kafka.WithSubscriberInitTopic(&sarama.TopicDetail{
			NumPartitions:     int32(cfg.SubscriberCfg.TopicDetails.NumPartitions),
			ReplicationFactor: int16(cfg.SubscriberCfg.TopicDetails.ReplicationFactor),
		}),
		kafka.WithSubscriberConsumerGroup("monitor"),
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

	monitorSubscriber, err := subscriber.NewMonitorSubscriber(
		&subscriber.Config{
			ProcessTopic: cfg.SubscriberCfg.ProcessTopic,
		},
		l,
		kafkaSubscriber,
		kafkaRouter,
		monitorPublisher,
	)
	if err != nil {
		l.Fatal("failed to create new monitor subscriber", map[string]interface{}{
			"error": err,
		})
	}

	monitorSubscriber.RegisterProcessHandler()

	go func() {
		if err := monitorSubscriber.Run(ctx); err != nil {
			l.Fatal("failed to run transaction subscriber", map[string]interface{}{
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
