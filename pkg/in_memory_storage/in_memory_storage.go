package in_memory_storage

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/pkg/collections"
	"applicationDesignTest/pkg/utils"
	"context"
	"fmt"
	"sync"
	"time"
)

var AvailabilitySeed = []model.Room{
	{"reddison", "lux", utils.Date(2024, 1, 1), 1},
	{"reddison", "lux", utils.Date(2024, 1, 2), 1},
	{"reddison", "lux", utils.Date(2024, 1, 3), 1},
	{"reddison", "lux", utils.Date(2024, 1, 4), 1},
	{"reddison", "lux", utils.Date(2024, 1, 5), 0},
}

type InMemoryStorage interface {
	OrderExists(ctx context.Context, order model.Order) bool
	CancelOrder(ctx context.Context, order model.Order) error
	CreatOrder(ctx context.Context, order model.Order) error

	CheckAvailability(_ context.Context, hotelID, roomID string, from, to time.Time) error
	IncrAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error
	DecrAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error
}

type inMemoryStorage struct {
	orders collections.ConcurrentMap[uint64, *model.Order]
	rooms  *sync.Map
}

func NewInMemoryStorage(rooms []model.Room) InMemoryStorage {
	roomsMap := sync.Map{}
	for _, room := range rooms {
		key := mapKey(room.HotelID, room.RoomID, room.Date)
		roomsMap.Store(key, room.Quota)
	}

	return &inMemoryStorage{
		orders: collections.NewConcurrentMap[uint64, *model.Order](10),
		rooms:  &roomsMap,
	}
}

func (m *inMemoryStorage) CreatOrder(_ context.Context, order model.Order) error {
	return m.orders.SetIfNotExists(order.ID, &order)
}

func (m *inMemoryStorage) OrderExists(_ context.Context, order model.Order) bool {
	return m.orders.Exists(order.ID)
}

func (m *inMemoryStorage) CancelOrder(_ context.Context, order model.Order) error {
	m.orders.Remove(order.ID)
	return nil
}

func (m *inMemoryStorage) CheckAvailability(_ context.Context, hotelID, roomID string, from, to time.Time) error {
	days := utils.DaysBetween(from, to)
	for _, day := range days {
		var avail int64 = 0

		key := mapKey(hotelID, roomID, day)
		v, ok := m.rooms.Load(key)
		if ok {
			avail = v.(int64)
		}

		if !ok || avail <= 0 {
			return fmt.Errorf("room key %s unavailable", key)
		}
	}

	return nil
}

func (m *inMemoryStorage) DecrAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error {
	days := utils.DaysBetween(from, to)
	for index, day := range days {
		key := mapKey(hotelID, roomID, day)
		err := m.decrAvailabilityCAS(key)
		if err != nil {
			// откатываем все изменения которые уже успели внести
			for i := 0; i < index; i++ {
				m.incrAvailabilityCAS(key)
			}

			return fmt.Errorf("room key %s unavailable", key)
		}
	}

	return nil
}

func (m *inMemoryStorage) IncrAvailability(ctx context.Context, hotelID, roomID string, from, to time.Time) error {
	days := utils.DaysBetween(from, to)
	for _, day := range days {
		key := mapKey(hotelID, roomID, day)
		m.incrAvailabilityCAS(key)
	}
	return nil
}

func mapKey(hotelID, roomID string, day time.Time) string {
	return fmt.Sprintf("%s:%s:%s", hotelID, roomID, utils.DateString(day))
}

func (m *inMemoryStorage) decrAvailabilityCAS(key string) error {
	done := false
	for !done {
		v, _ := m.rooms.Load(key)
		val := v.(int64)
		if val == 0 {
			return fmt.Errorf("room key %s unavailable", key)
		}

		done = m.rooms.CompareAndSwap(key, val, val-1)
	}

	return nil
}

func (m *inMemoryStorage) incrAvailabilityCAS(key string) {
	done := false
	for !done {
		v, _ := m.rooms.Load(key)
		val := v.(int64)
		done = m.rooms.CompareAndSwap(key, val, val-1)
	}
}
