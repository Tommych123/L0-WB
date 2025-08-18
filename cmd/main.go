package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Tommych123/L0-WB/internal/config"
	httphandler "github.com/Tommych123/L0-WB/internal/http"
	"github.com/Tommych123/L0-WB/internal/kafka"
	"github.com/Tommych123/L0-WB/internal/repository"
	"github.com/Tommych123/L0-WB/internal/service"
)

func main() {
	// Загружаем конфиг
	cfg := config.Load()

	// Подключение к БД
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Репозиторий
	repo := repository.NewPostgresOrderRepository(db)

	// Сервис заказов
	orderService := service.NewOrderService(repo)

	// Загрузка кеша из БД при старте
	if err := orderService.LoadCache(); err != nil {
		log.Fatalf("failed to load cache: %v", err)
	}

	// Kafka consumer
	consumer := kafka.NewConsumer(cfg.KafkaBroker, "orders-topic", "orders-group", orderService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск Kafka consumer
	go func() {
		if err := consumer.Run(ctx); err != nil {
			log.Printf("Kafka consumer stopped: %v", err)
		}
	}()

	// HTTP сервер
	router := httphandler.NewRouter(orderService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServicePort),
		Handler: router,
	}

	go func() {
		log.Printf("HTTP server listening on :%d", cfg.ServicePort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Завершаем работу HTTP сервера
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
