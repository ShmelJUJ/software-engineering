package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/subscriber"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/generated/restapi"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/generated/restapi/operations"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/repository"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/usecase"

	apiTransaction "github.com/ShmelJUJ/software-engineering/transaction/internal/generated/restapi/operations/transaction"

	monitor_client "github.com/ShmelJUJ/software-engineering/pkg/monitor_client/client"
	"github.com/go-openapi/loads"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/ShmelJUJ/software-engineering/pkg/kafka"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ShmelJUJ/software-engineering/pkg/postgres"
	"github.com/ShmelJUJ/software-engineering/pkg/redis"
	"github.com/ShmelJUJ/software-engineering/transaction/config"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/api/handler"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/api/middleware"
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

	l, err := logger.NewLogrusLogger(cfg.LoggerCfg.Level)
	if err != nil {
		log.Fatalln("failed to create new logger: ", err)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		l.Fatal("failed to get swagger spec", map[string]interface{}{
			"error": err,
		})
	}

	pg, err := postgres.New(ctx, cfg.PostgresCfg.URL, postgres.WithLogger(l))
	if err != nil {
		l.Fatal("failed to create a new postgre", map[string]interface{}{
			"error": err,
		})
	}
	defer pg.Close()

	r, err := redis.New(ctx, &redis.Config{
		ConnURL: cfg.RedisCfg.URL,
	})
	if err != nil {
		l.Fatal("failed to create a new redis", map[string]interface{}{
			"error": err,
		})
	}
	defer r.Close()

	monitorClientCfg := monitor_client.DefaultTransportConfig().
		WithBasePath("api/v1").
		WithHost("host.docker.internal:8080").
		WithSchemes([]string{"http"})

	transport := httptransport.New(monitorClientCfg.Host, monitorClientCfg.BasePath, monitorClientCfg.Schemes)

	monitorClient := monitor_client.New(transport, strfmt.Default)

	kafkaPublisher, err := kafka.NewPublisher(cfg.PublisherCfg.Brokers)
	if err != nil {
		l.Fatal("failed to create kafka publisher", map[string]interface{}{
			"error": err,
		})
	}

	transactionPublisher, err := publisher.NewTransactionPublisher(
		&publisher.Config{
			ProcessedTransactionTopic: cfg.PublisherCfg.ProcessedTransactionTopic,
		},
		l,
		kafkaPublisher,
	)
	if err != nil {
		l.Fatal("failed to create a transaction publisher", map[string]interface{}{
			"error": err,
		})
	}

	transactionRepo := repository.NewTransactionRepo(pg, l)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, transactionPublisher, l)
	transactionHandler := handler.NewTransactionHandler(
		transactionUsecase,
		l,
		monitorClient.Monitor,
		r,
	)

	middlewareManager, err := middleware.NewMiddlewareManager(&middleware.Config{
		IdempotencyCfg: &middleware.IdempotencyConfig{
			Name:      cfg.MiddlewareCfg.IdempotenctCfg.Name,
			HeaderKey: cfg.MiddlewareCfg.IdempotenctCfg.HeaderKey,
		},
	}, l, r)
	if err != nil {
		l.Fatal("failed to create middleware manager", map[string]interface{}{
			"error": err,
		})
	}

	api := operations.NewTransactionAPI(swaggerSpec)

	api.BearerAuth = transactionHandler.VerifyAuthToken
	api.TransactionAcceptTransactionHandler = apiTransaction.AcceptTransactionHandlerFunc(transactionHandler.AcceptTransactionHandler)
	api.TransactionCancelTransactionHandler = apiTransaction.CancelTransactionHandlerFunc(transactionHandler.CancelTransactionHandler)
	api.TransactionCreateTransactionHandler = apiTransaction.CreateTransactionHandlerFunc(transactionHandler.CreateTransactionHandler)
	api.TransactionEditTransactionHandler = apiTransaction.EditTransactionHandlerFunc(transactionHandler.EditTransactionHandler)
	api.TransactionRetrieveTransactionHandler = apiTransaction.RetrieveTransactionHandlerFunc(transactionHandler.RetrieveTransactionHandler)
	api.TransactionRetrieveTransactionStatusHandler = apiTransaction.RetrieveTransactionStatusHandlerFunc(transactionHandler.RetrieveTransactionStatusHandler)
	api.TransactionLoginHandler = apiTransaction.LoginHandlerFunc(transactionHandler.LoginHandler)

	middlewareManager.AddIdempotenceMiddleware()
	middlewareManager.SetupGlobalMiddleware(swaggerSpec, api)

	server := restapi.NewServer(api)
	server.EnabledListeners = []string{"unix", "http"}

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

	defer func() {
		if err := server.Shutdown(); err != nil {
			l.Error("failed to shutdown server", map[string]interface{}{
				"error": err,
			})
		}
	}()

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
		kafka.WithSubscriberConsumerGroup("transaction"),
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

	transactionSub, err := subscriber.NewTransactionSubscriber(
		&subscriber.Config{
			FailedTransactionTopic:    cfg.SubscriberCfg.FailedTransactionTopic,
			SucceededTransactionTopic: cfg.SubscriberCfg.SucceededTransactionTopic,
		},
		l,
		kafkaSubscriber,
		kafkaRouter,
		transactionRepo,
	)
	if err != nil {
		l.Fatal("failed to create new transaction subscriber", map[string]interface{}{
			"error": err,
		})
	}

	transactionSub.RegisterFailedTransactionHandler()
	transactionSub.RegisterSucceededTransactionHandler()

	go func() {
		if err := transactionSub.Run(ctx); err != nil {
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
