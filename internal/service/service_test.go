package service

import (
	"errors"
	"testing"

	"github.com/Tommych123/L0-WB/internal/domain"
)

// --- mockRepo ---
type mockRepo struct {
	saveFunc   func(order *domain.Order) error
	getFunc    func(id string) (*domain.Order, error)
	getAllFunc func() ([]*domain.Order, error)
}

func (m *mockRepo) Save(order *domain.Order) error {
	if m.saveFunc != nil {
		return m.saveFunc(order)
	}
	return nil
}
func (m *mockRepo) Get(id string) (*domain.Order, error) {
	if m.getFunc != nil {
		return m.getFunc(id)
	}
	return nil, nil
}
func (m *mockRepo) GetAll() ([]*domain.Order, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return nil, nil
}

// --- Тесты ---

// SaveOrder успешный
func TestSaveOrder_Success(t *testing.T) {
	mock := &mockRepo{}
	s := NewOrderService(mock)

	order := &domain.Order{OrderUID: "test123"}
	if err := s.SaveOrder(order); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// SaveOrder с ошибкой репозитория
func TestSaveOrder_RepoError(t *testing.T) {
	mock := &mockRepo{
		saveFunc: func(order *domain.Order) error {
			return errors.New("repo error")
		},
	}
	s := NewOrderService(mock)

	order := &domain.Order{OrderUID: "test123"}
	if err := s.SaveOrder(order); err == nil {
		t.Fatal("expected error, got nil")
	}
}

// GetOrder успешный
func TestGetOrder_Success(t *testing.T) {
	mock := &mockRepo{
		getFunc: func(id string) (*domain.Order, error) {
			return &domain.Order{OrderUID: id}, nil
		},
	}
	s := NewOrderService(mock)

	got, err := s.GetOrder("order1")
	if err != nil || got.OrderUID != "order1" {
		t.Errorf("unexpected result: %v, %v", got, err)
	}
}

// GetOrder не найден
func TestGetOrder_NotFound(t *testing.T) {
	mock := &mockRepo{
		getFunc: func(id string) (*domain.Order, error) {
			return nil, errors.New("not found")
		},
	}
	s := NewOrderService(mock)

	got, err := s.GetOrder("order1")
	if got != nil || err == nil {
		t.Errorf("expected nil and error, got: %v, %v", got, err)
	}
}

// SaveOrder кладет в кеш
func TestSaveOrder_Cache(t *testing.T) {
	mock := &mockRepo{}
	s := NewOrderService(mock)

	order := &domain.Order{OrderUID: "cache1"}
	if err := s.SaveOrder(order); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := s.GetOrder("cache1")
	if got == nil || got.OrderUID != "cache1" {
		t.Errorf("order not found in cache after SaveOrder")
	}
}

// LoadCache загружает все заказы в кеш
func TestLoadCache(t *testing.T) {
	mockOrders := []*domain.Order{
		{OrderUID: "order1"},
		{OrderUID: "order2"},
	}
	mock := &mockRepo{
		getAllFunc: func() ([]*domain.Order, error) {
			return mockOrders, nil
		},
	}
	s := NewOrderService(mock)

	if err := s.LoadCache(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, o := range mockOrders {
		if got, _ := s.GetOrder(o.OrderUID); got == nil {
			t.Errorf("order %s not loaded into cache", o.OrderUID)
		}
	}
}

// GetOrder берет из кеша
func TestGetOrder_FromCache(t *testing.T) {
	mock := &mockRepo{}
	s := NewOrderService(mock)

	order := &domain.Order{OrderUID: "cached"}
	s.cache["cached"] = order // вручную положили в кеш

	got, err := s.GetOrder("cached")
	if err != nil || got.OrderUID != "cached" {
		t.Errorf("expected to get order from cache, got: %v, %v", got, err)
	}
}
