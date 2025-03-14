package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/Levan1e/url-shortener-service/internal/api"
	"github.com/Levan1e/url-shortener-service/internal/repository/memory"
	"github.com/Levan1e/url-shortener-service/internal/repository/postgres"
	"github.com/Levan1e/url-shortener-service/internal/service"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	storageType := flag.String("storage", "memory", "Тип хранилища: memory или postgres")
	flag.Parse()

	var store service.Storage

	if *storageType == "postgres" {
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			logrus.Fatal("Для postgres хранилища необходимо задать переменную окружения DATABASE_URL")
		}
		pStore, err := postgres.NewStorage(connStr)
		if err != nil {
			logrus.Fatalf("Ошибка инициализации PostgreSQL хранилища: %v", err)
		}
		store = pStore
	} else {
		store = memory.NewStorage()
	}
	shortenerService := service.NewShortenerService(store)
	handler := api.NewHandler(shortenerService)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		logrus.WithFields(logrus.Fields{
			"port":    port,
			"storage": *storageType,
		}).Info("Сервис запущен")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logrus.Info("Получен сигнал завершения, начинаем корректное завершение работы...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Ошибка при завершении работы сервера: %v", err)
	}

	logrus.Info("Сервер завершил работу")
}
