package repository

import "github.com/Tommych123/L0-WB/internal/domain"

// Интерфейс для работы с БД
type OrderRepository interface {
	Save(order *domain.Order) error
	Get(orderUID string) (*domain.Order, error)
	GetAll() ([]*domain.Order, error)
}
