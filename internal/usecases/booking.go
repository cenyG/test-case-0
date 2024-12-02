package usecases

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/repo"
	"context"
	"errors"
	"fmt"
)

type RetryBookingQueue interface {
	Enqueue(order model.Order)
}

type BookingUseCase interface {
	BookRoom(ctx context.Context, order model.Order) error
}

type bookingUseCase struct {
	ordersRepo repo.OrdersRepo
	roomsRepo  repo.RoomsRepo
	retryQueue RetryBookingQueue
}

func NewBookingUseCase(ordersRepo repo.OrdersRepo, booksRepo repo.RoomsRepo, retryQueue RetryBookingQueue) BookingUseCase {
	return &bookingUseCase{
		ordersRepo: ordersRepo,
		roomsRepo:  booksRepo,
		retryQueue: retryQueue,
	}
}

func (b *bookingUseCase) BookRoom(ctx context.Context, order model.Order) error {
	exists := b.ordersRepo.OrderExists(ctx, order)
	if exists {
		// если уже создан
		return nil
	}

	err := b.roomsRepo.CheckAvailability(ctx, order.HotelID, order.RoomID, order.From, order.To)
	if err != nil {
		return fmt.Errorf("b.roomsRepo.CheckAvailability(order=%d): %v", order.ID, err)
	}

	err = b.ordersRepo.CreatOrder(ctx, order)
	if err != nil {
		if errors.Is(err, repo.ErrOrderAlreadyCreated) {
			return nil
		}
		// пытаемся добить order
		b.retryQueue.Enqueue(order)

		return fmt.Errorf("b.ordersRepo.CreatOrder() err: %v, will retry", order)
	}

	err = b.roomsRepo.DecrAvailability(ctx, order.HotelID, order.RoomID, order.From, order.To)
	if err != nil {
		cancelErr := b.ordersRepo.CancelOrder(ctx, order)
		if cancelErr != nil {
			return fmt.Errorf("b.ordersRepo.CancelOrder(order=%d): %v", order.ID, err)
		}
		return fmt.Errorf("b.roomsRepo.DecrAvailability(order=%d): %v", order.ID, err)
	}

	return err
}
