package in_memory

import (
	"applicationDesignTest/internal/repo"
	"applicationDesignTest/pkg/in_memory_storage"
	"context"
	"time"
)

type InMemoryRoomsRepo struct {
	storage in_memory_storage.InMemoryStorage
}

func NewInMemoryRoomsRepo(storage in_memory_storage.InMemoryStorage) repo.RoomsRepo {
	return &InMemoryRoomsRepo{
		storage: storage,
	}
}

func (o *InMemoryRoomsRepo) IncrAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error {
	return o.storage.IncrAvailability(ctx, hotelID, roomID, from, to)
}

func (o *InMemoryRoomsRepo) DecrAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error {
	return o.storage.DecrAvailability(ctx, hotelID, roomID, from, to)
}

func (o *InMemoryRoomsRepo) CheckAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error {
	return o.storage.CheckAvailability(ctx, hotelID, roomID, from, to)
}
