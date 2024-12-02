package repo

import (
	"applicationDesignTest/internal/model"
	"context"
	"errors"
	"time"
)

var (
	ErrOrderAlreadyCreated = errors.New("order already created")
)

type OrdersRepo interface {
	CreatOrder(ctx context.Context, order model.Order) error
	CancelOrder(ctx context.Context, order model.Order) error
	OrderExists(ctx context.Context, order model.Order) bool
}

type RoomsRepo interface {
	IncrAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error
	DecrAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error
	CheckAvailability(_ context.Context, hotelID, roomID string, from, to time.Time) error
}
