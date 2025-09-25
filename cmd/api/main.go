package main

import (
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pirasl/payment-service/internal/data"
)

type service struct {
	config         *serviceConfig
	logger         *slog.Logger
	models         *data.Models
	rabbitmqClient *rabbitMQClient

	stripeClient *stripeConfig
	workerPool   *workerPool

	wg sync.WaitGroup
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(slog.String("service", "[PAYMENT-STRIPE]"))

	logger.Info("starting...")

	if getOptionalStringEnv("APP_ENV", "development") == "development" {
		logger.Info("[DEV MODE]: loading env file")
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file, using OS environment variables")
		}
	}

	logger.Info("reaching for rabbitmq...")
	rabbitmq, err := newRabbitMQClient()
	if err != nil {
		logger.Error("could not reach rabbitmq. exiting...", "err: ", err)
		os.Exit(1)
	}
	logger.Info("connected to rabbitmq.")

	logger.Info("loading stripe config...")

	stripeClient, err := newStripeConfig()
	if err != nil {
		logger.Error("cannot load stripe config. exiting...", "err:", err)
		os.Exit(1)
	}

	logger.Info("connecting to transactions db...")

	db, err := openDB()
	if err != nil {
		logger.Error("failed to open postgres db. exiting...", "err:", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		logger.Error("failed to run migrations. exiting...", "err:", err)
		os.Exit(1)
	}

	logger.Info("connected to stripe payment db.")

	serviceConfig, err := newServiceConfig()
	if err != nil {
		logger.Error("cannot load service config", "err:", err)
		os.Exit(1)
	}

	wp, err := newWorkerPool(5, rabbitmq)
	if err != nil {
		logger.Error("unable to create worker pool. exiting...", "err:", err)
	}

	s := &service{
		logger:         logger,
		rabbitmqClient: rabbitmq,
		models:         data.NewModels(db),
		stripeClient:   stripeClient,
		config:         serviceConfig,
		workerPool:     wp,
	}

	logger.Info("service config loaded")

	go s.gRPCListen()

	logger.Info("stripe payment service up and running", "port", serviceConfig.gRPCPort)

	if err := s.serve(); err != nil {
		logger.Error("error while starting server", "err", err)
		os.Exit(1)
	}

}
