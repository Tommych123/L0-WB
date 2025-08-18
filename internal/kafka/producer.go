package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Tommych123/L0-WB/internal/domain"
	"github.com/segmentio/kafka-go"
)

// Отвечает за отправку сообщений в Kafka
type Producer struct {
	writer *kafka.Writer
}

// Создает нового продюсера
func NewProducer(broker, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// Отправляет заказ в Kafka в формате JSON
func (p *Producer) SendOrder(ctx context.Context, order domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return err
	}

	log.Printf("Order sent to Kafka: %s", order.OrderUID)
	return nil
}

// Закрывает соединение с Kafka
func (p *Producer) Close() error {
	return p.writer.Close()
}
