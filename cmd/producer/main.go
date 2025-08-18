package main

import (
	"context"
	"log"
	"time"

	"github.com/Tommych123/L0-WB/internal/domain"
	"github.com/Tommych123/L0-WB/internal/kafka"
)

// Пример приходящего заказа для теста обработки сообщений
func main() {
	producer := kafka.NewProducer("localhost:9092", "orders-topic")
	defer producer.Close()
	// Заказ
	order := domain.Order{
		OrderUID:    "test12345",
		TrackNumber: "WBTESTTRACK",
		Entry:       "WBIL",
		Delivery: domain.Delivery{
			Name:    "Test User",
			Phone:   "+123456789",
			Zip:     "123456",
			City:    "Test City",
			Address: "Test Street 1",
			Region:  "Test Region",
			Email:   "test@example.com",
		},
		Payment: domain.Payment{
			Transaction:  "test12345",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       100,
			PaymentDt:    time.Now().Unix(),
			Bank:         "Test Bank",
			DeliveryCost: 10,
			GoodsTotal:   90,
			CustomFee:    0,
		},
		Items: []domain.Item{
			{
				ChrtID:      1,
				TrackNumber: "WBTESTTRACK",
				Price:       90,
				Name:        "Test Item",
				Sale:        0,
				Size:        "M",
				TotalPrice:  90,
				NmID:        1001,
				Brand:       "Test Brand",
				Status:      1,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "cust1",
		DeliveryService:   "meest",
		ShardKey:          "1",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	// Отправка его
	err := producer.SendOrder(context.Background(), order)
	if err != nil {
		log.Fatalf("failed to send order: %v", err)
	}
}
