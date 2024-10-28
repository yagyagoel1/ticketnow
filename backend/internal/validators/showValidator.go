package validators

import "time"

type CreateShow struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	ShowTiming  time.Time `json:"showTiming" validate:"required"`
	TicketTypes []struct {
		Name  string  `json:"name" validate:"required"`
		Price float64 `json:"price" validate:"required"`
		Count uint    `json:"count" validate:"required"`
	} `json:"ticketTypes" validate:"required"`
}
