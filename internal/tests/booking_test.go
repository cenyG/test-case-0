package tests

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/repo/in_memory"
	"applicationDesignTest/internal/services"
	"applicationDesignTest/internal/usecases"
	"applicationDesignTest/pkg/in_memory_storage"
	"applicationDesignTest/pkg/utils"
	"context"
	"sync"
	"testing"
	"time"
)

func Test_bookingUseCase_BookRoom(t *testing.T) {

	type testCase struct {
		name           string
		availabilities []model.Room
		orders         []model.Order
		wantErr        bool
	}

	ctx := context.Background()

	tests := []testCase{
		{
			name:   "test parallel",
			orders: generateOrders(100, "reddison", "lux", utils.Date(2024, 1, 1), utils.Date(2024, 1, 4)),
			availabilities: []model.Room{
				{"reddison", "lux", utils.Date(2024, 1, 1), 100},
				{"reddison", "lux", utils.Date(2024, 1, 2), 100},
				{"reddison", "lux", utils.Date(2024, 1, 3), 100},
				{"reddison", "lux", utils.Date(2024, 1, 4), 100},
				{"reddison", "lux", utils.Date(2024, 1, 5), 0},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Storage
			storage := in_memory_storage.NewInMemoryStorage(tt.availabilities)

			// Repos
			ordersRepo := in_memory.NewInMemoryOrdersRepo(storage)
			roomsRepo := in_memory.NewInMemoryRoomsRepo(storage)

			//UseCases and Services
			retryBooking := usecases.NewRetryBookingUseCase(ordersRepo, roomsRepo)
			retryWorkerService := services.NewRetryCreateOrderWorker(retryBooking, services.RetryConfig{
				RetryCount:  3,
				Delay:       100 * time.Millisecond,
				QueueLength: 1000,
				TTL:         2 * time.Minute,
			})
			booking := usecases.NewBookingUseCase(ordersRepo, roomsRepo, &retryWorkerService)

			wg := sync.WaitGroup{}
			wg.Add(len(tt.orders))

			for _, v := range tt.orders {
				go func(order model.Order) {
					defer wg.Done()

					err := booking.BookRoom(ctx, order)
					if err != nil {
						t.Errorf("error: %v", err)
						return
					}
				}(v)
			}
			wg.Wait()

			err := roomsRepo.CheckAvailability(ctx, tt.orders[0].HotelID, tt.orders[0].RoomID, tt.orders[0].From, tt.orders[0].To)
			if err == nil {
				t.Fatalf("must be no empty rooms")
			}

		})
	}
}

func generateOrders(count int, hotelID, roomID string, from, to time.Time) []model.Order {
	orders := make([]model.Order, count)
	for i, _ := range orders {
		orders[i] = model.Order{
			ID:        uint64(i),
			HotelID:   hotelID,
			RoomID:    roomID,
			UserEmail: "",
			From:      from,
			To:        to,
		}
	}

	return orders
}
