package usecases

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/repo"
	"context"
	"errors"
	"fmt"
)

type RetryBookingUseCase interface {
	RetryBooking(ctx context.Context, order model.Order) error
}

type retryBookingUseCase struct {
	ordersRepo repo.OrdersRepo
	roomsRepo  repo.RoomsRepo
}

func NewRetryBookingUseCase(ordersRepo repo.OrdersRepo, roomsRepo repo.RoomsRepo) RetryBookingUseCase {
	return &retryBookingUseCase{
		ordersRepo: ordersRepo,
		roomsRepo:  roomsRepo,
	}
}

func (r *retryBookingUseCase) RetryBooking(ctx context.Context, order model.Order) error {
	err := r.ordersRepo.CreatOrder(ctx, order)
	if err != nil {
		if errors.Is(err, repo.ErrOrderAlreadyCreated) {
			return nil
		}
		return fmt.Errorf("r.ordersRepo.CreatOrder(%d): %v", order.ID, err)
	}

	err = r.roomsRepo.DecrAvailability(ctx, order.HotelID, order.RoomID, order.From, order.To)
	if err != nil {
		cancelErr := r.ordersRepo.CancelOrder(ctx, order)
		if cancelErr != nil {
			return fmt.Errorf("r.ordersRepo.CancelOrder(%d): %v", order.ID, err)
		}
		return fmt.Errorf("r.roomsRepo.IncrAvailability(%d): %v", order.ID, err)
	}

	return nil
}
