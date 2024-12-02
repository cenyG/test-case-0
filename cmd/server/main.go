package main

import (
	"applicationDesignTest/pkg/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"applicationDesignTest/internal/repo/in_memory"
	"applicationDesignTest/internal/services"
	"applicationDesignTest/internal/transport/http"
	"applicationDesignTest/internal/usecases"
	"applicationDesignTest/pkg/httpserver"
	"applicationDesignTest/pkg/in_memory_storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Storage
	storage := in_memory_storage.NewInMemoryStorage(in_memory_storage.AvailabilitySeed)

	// Repos
	ordersRepo := in_memory.NewInMemoryOrdersRepo(storage)
	roomsRepo := in_memory.NewInMemoryRoomsRepo(storage)

	//UseCases and Services
	retryBookingUseCase := usecases.NewRetryBookingUseCase(ordersRepo, roomsRepo)
	retryWorkerService := services.NewRetryCreateOrderWorker(retryBookingUseCase, services.RetryConfig{
		RetryCount:  3,
		Delay:       100 * time.Millisecond,
		QueueLength: 1000,
		TTL:         2 * time.Minute,
	})
	bookingUseCase := usecases.NewBookingUseCase(ordersRepo, roomsRepo, &retryWorkerService)

	useCaseContainer := usecases.UseCaseContainer{
		BookingUseCase:      bookingUseCase,
		RetryBookingUseCase: retryBookingUseCase,
	}
	// HTTP Server
	handler := gin.New()
	http.NewRouter(handler, useCaseContainer)
	httpServer := httpserver.New(handler, httpserver.Port("8080"))

	// Start server
	utils.Go(ctx, func(ctx context.Context) {
		httpServer.Start()
	})

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Error(fmt.Sprintf("[main] signal %s", s))
		cancel()
	case err := <-httpServer.Notify():
		slog.Error(fmt.Sprintf("[main] httpServer.Notify: %v", err))
		cancel()
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		slog.Error(fmt.Sprintf("[main] httpServer.Shutdown: %v", err))
	}
}
