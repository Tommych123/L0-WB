package service

import (
	"log"
	"sync"

	"github.com/Tommych123/L0-WB/internal/domain"
	"github.com/Tommych123/L0-WB/internal/repository"
)

// Структура для сервиса
type OrderService struct {
	repo  repository.OrderRepository
	cache map[string]*domain.Order
	mu    sync.RWMutex
}

// Создание нового сервиса
func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: make(map[string]*domain.Order),
	}
}

// Сохраняет заказ в БД и кэш
func (s *OrderService) SaveOrder(order *domain.Order) error {
	if err := s.repo.Save(order); err != nil {
		return err
	}

	s.mu.Lock()
	s.cache[order.OrderUID] = order
	s.mu.Unlock()

	return nil
}

// Возвращает заказ из кэша или из БД
func (s *OrderService) GetOrder(id string) (*domain.Order, error) {
	s.mu.RLock()
	if order, exists := s.cache[id]; exists {
		s.mu.RUnlock()
		return order, nil
	}
	s.mu.RUnlock()

	order, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if order != nil {
		s.mu.Lock()
		s.cache[id] = order
		s.mu.Unlock()
	}

	return order, nil
}

// Заполняет кэш из БД при старте сервиса
func (s *OrderService) LoadCache() error {
	orders, err := s.repo.GetAll()
	if err != nil {
		return err
	}

	s.mu.Lock()
	for _, order := range orders {
		s.cache[order.OrderUID] = order
	}
	s.mu.Unlock()

	log.Printf("Cache loaded with %d orders", len(orders))
	return nil
}
