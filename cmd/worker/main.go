package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/worker"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/di"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	awsConfig := config.LoadAWSConfig(cfg)
	config.InstantiateSQSClient(awsConfig)

	db, err := di.NewDatabaseConnectionPool(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot initialize zap logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck // It is not necessary to check for errors at this moment.

	zap.ReplaceGlobals(logger.With(zap.String("app", cfg.AppName), zap.String("env", cfg.AppEnv)))

	// Datastore.
	ds := datastore.NewDatastore(db)

	// Use cases.
	updateOrderUseCase := usecase.NewUpdateOrderUseCase(ds)

	// Consumers.
	updateOrderStatusConsumer := worker.NewUpdateOrderStatusConnsumer(updateOrderUseCase)

	chnProcessingMessages := make(chan worker.ProcessingMessage, awsConfig.SQSConfig.NumWorkers)
	ctx, cancel := context.WithCancel(context.Background())

	var wgConsumer, wgProcessing sync.WaitGroup

	for range awsConfig.SQSConfig.NumWorkers {
		wgProcessing.Add(1)
		go startWorker(ctx, chnProcessingMessages, &wgProcessing, *awsConfig, updateOrderStatusConsumer)
	}

	for _, sqsQueue := range awsConfig.SQSConfig.Queues {
		wgConsumer.Add(1)
		go worker.Consumer(ctx, *awsConfig, sqsQueue, chnProcessingMessages, &wgConsumer)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		select {
		case signalReceived := <-signalCh:
			zap.L().Info("Received signal", zap.String("signal", signalReceived.String()))
			zap.L().Info("Shutting down gracefully...")

			cancel()
		case <-ctx.Done():
			zap.L().Info("Context canceled, shutting down gracefully...")
		}

		wgConsumer.Wait()

		zap.L().Info("All consumers have finished.")

		close(chnProcessingMessages)
	}()

	zap.L().Info("Starting consumers...")

	wgProcessing.Wait()

	zap.L().Info("All workers have finished.")
}

func startWorker(
	ctx context.Context,
	chnProcessingMessages <-chan worker.ProcessingMessage,
	wg *sync.WaitGroup,
	awsConfig config.AWSConfig,
	updateOrderStatusConsumer worker.Processor,
) {
	defer wg.Done()

	for processingMessage := range chnProcessingMessages {
		zap.L().Info("Processing message from queue", zap.String("queue", processingMessage.QueueName))

		var processingError error

		func() {
			defer func() {
				if r := recover(); r != nil {
					zap.L().Error("Recovered from panic while processing message", zap.Any("error", r))
				}
			}()

			if processingMessage.QueueName == "payment_events" {
				processingError = updateOrderStatusConsumer.Process(ctx, processingMessage)
			}
		}()

		if processingError == nil {
			deleteMessageErr := worker.DeleteMessage(awsConfig, processingMessage)
			if deleteMessageErr != nil {
				zap.L().Error(
					"Error deleting message from queue",
					zap.String("queue", processingMessage.QueueName),
					zap.Error(deleteMessageErr),
				)
			}
		}
	}
}
