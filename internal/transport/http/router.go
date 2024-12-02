package http

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/transport/http/handlers"
	"applicationDesignTest/internal/transport/http/middleware"
	"applicationDesignTest/internal/usecases"
	"github.com/gin-gonic/gin"
)

// NewRouter - setup Gin router
func NewRouter(handler *gin.Engine, container usecases.UseCaseContainer) {
	handler.Use(gin.Logger())
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Handlers
	ordersHandler := handlers.NewOrdersHandler(container.BookingUseCase)

	// Routers
	group := handler.Group("/orders")
	group.POST("/", middleware.JSONDeserializerMiddleware(&model.Order{}), ordersHandler.Handle)
}
