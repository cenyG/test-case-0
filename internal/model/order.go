package model

import "time"

type Order struct {
	ID        uint64    `json:"id"` // prefer uuid
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}
