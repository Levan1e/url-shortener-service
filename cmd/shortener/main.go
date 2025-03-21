package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Levan1e/url-shortener-service/internal/api"
	v1 "github.com/Levan1e/url-shortener-service/internal/api/v1"
	"github.com/Levan1e/url-shortener-service/internal/config"
	"github.com/Levan1e/url-shortener-service/internal/repository/memory"
	"github.com/Levan1e/url-shortener-service/internal/repository/postgres"
	"github.com/Levan1e/url-shortener-service/internal/service"
	"github.com/Levan1e/url-shortener-service/pkg/logger"
	postgres_helpers "github.com/Levan1e/url-shortener-service/pkg/postgres"
)

func main() {
	config, err := config.ParseConfig("internal/config/config.yaml")
	if err != nil {
		logger.Fatalf("Failed to parse config: %v", err)
	}

	storageType := flag.String("storage", "memory", "Тип хранилища: memory или postgres")
	flag.Parse()

	if storageType == nil {
		logger.Fatal("Storage must not be empty")
	}

	ctx := context.Background()

	var memStorage *memory.MemoryStorage
	var storage service.Storage

	switch *storageType {
	case "memory":
		memStorage = memory.NewStorage()
		storage = memStorage
	case "postgres":
		pool, err := postgres_helpers.NewPostgresPool(ctx, config.Postgres)
		if err != nil {
			logger.Fatalf("Failed to create postgres pool: %v", err)
		}
		postgres_helpers.Migrate(pool, config.MigrationsDir)
		storage = postgres.NewStorage(pool)
	default:
		logger.Fatalf("Storage must be memory or postgres")
	}

	shortenerService := service.NewShortenerService(storage)
	v1Handler := v1.NewHandler(shortenerService)
	server := api.NewServer(config.Server, v1Handler)

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logrus.Fatalf("Failed to stop server: %v", err)
	}

	if memStorage != nil {
		if err := memStorage.SaveToFileOnShutdown("storage.json"); err != nil {
			logger.Errorf("Failed to save memory storage to file: %v", err)
		} else {
			logger.Info("Memory storage saved to storage.json")
		}
	}
}
