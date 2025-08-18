package http

import (
	"github.com/Tommych123/L0-WB/internal/service"
	"github.com/gorilla/mux"
)

func NewRouter(orderService *service.OrderService) *mux.Router {
	r := mux.NewRouter()

	h := NewHandler(orderService)

	// Эндпоинт для выдачи заказа по ID
	r.HandleFunc("/orders/{id}", h.GetOrderByID).Methods("GET")

	// Веб страница
	r.HandleFunc("/", h.ServeWebUI).Methods("GET")

	return r
}
