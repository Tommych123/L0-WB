package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/Tommych123/L0-WB/internal/domain"
	"github.com/Tommych123/L0-WB/internal/service"
	"github.com/segmentio/kafka-go"
)

// Consumer читает сообщения из Kafka
type Consumer struct {
	reader       *kafka.Reader
	orderService *service.OrderService
}

// Создание нового Consumer
func NewConsumer(broker, topic, groupID string, orderService *service.OrderService) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			GroupID: groupID,
			Topic:   topic,
		}),
		orderService: orderService,
	}
}

// Проверка корректности заказа
func isValidOrder(order *domain.Order) bool {
	if order.OrderUID == "" {
		return false
	}
	if order.Payment.Amount < 0 || order.Payment.Currency == "" {
		return false
	}
	if order.Delivery.Email == "" || !strings.Contains(order.Delivery.Email, "@") {
		return false
	}
	if len(order.Items) == 0 {
		return false
	}

	for _, item := range order.Items {
		if item.Name == "" || item.Price < 0 || item.TotalPrice < 0 {
			return false
		}
	}

	return true
}

// Чтение сообщений из Kafka
func (c *Consumer) Run(ctx context.Context) error {
	for {
		// Получаем сообщение
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("error fetching message: %v", err)
			continue
		}

		// Парсим JSON
		var order domain.Order
		err = json.Unmarshal(m.Value, &order)
		if err != nil {
			log.Printf("invalid message format: %v", err)
			if commitErr := c.reader.CommitMessages(ctx, m); commitErr != nil {
				log.Printf("error committing invalid message: %v", commitErr)
			}
			continue
		}

		// Валидация заказа
		if !isValidOrder(&order) {
			log.Printf("invalid order data: %v", order.OrderUID)
			if commitErr := c.reader.CommitMessages(ctx, m); commitErr != nil {
				log.Printf("error committing invalid order: %v", commitErr)
			}
			continue
		}

		// Сохраняем заказ через сервис (БД + кэш)
		err = c.orderService.SaveOrder(&order)
		if err != nil {
			log.Printf("error saving order: %v", err)
			continue
		}

		// Подтверждаем успешную обработку
		if commitErr := c.reader.CommitMessages(ctx, m); commitErr != nil {
			log.Printf("error committing message: %v", commitErr)
		}

		log.Printf("order saved: %s", order.OrderUID)
	}
}
