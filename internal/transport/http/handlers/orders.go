package handlers

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ordersHandler struct {
	useCase usecases.BookingUseCase
}

func NewOrdersHandler(useCase usecases.BookingUseCase) Handler {
	return &ordersHandler{useCase}
}

// Handle - upload file handler
func (h *ordersHandler) Handle(c *gin.Context) {
	// Извлекаем десериализованные данные из контекста
	jsonData, exists := c.Get("jsonData")
	if !exists {
		failResponse(c, http.StatusInternalServerError, "error while JSON parsing")
		return
	}

	order, ok := jsonData.(*model.Order)
	if !ok {
		failResponse(c, http.StatusInternalServerError, "type mismatch for order")
		return
	}

	err := h.useCase.BookRoom(c, *order)
	if err != nil {
		failResponse(c, http.StatusInternalServerError, fmt.Sprintf("cant create order: %v", err))
		return
	}

	successResponse(c, order.ID)
}
