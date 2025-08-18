package http

import (
	"encoding/json"
	"net/http"

	"github.com/Tommych123/L0-WB/internal/service"
	"github.com/gorilla/mux"
)

// Структура handler
type Handler struct {
	orderService *service.OrderService
}

// Создание нового handler
func NewHandler(orderService *service.OrderService) *Handler {
	return &Handler{orderService: orderService}
}

// Получение заказа по ID
func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		http.Error(w, "Ошибка сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if order == nil {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// Отдает статическую HTML-страницу
func (h *Handler) ServeWebUI(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}
