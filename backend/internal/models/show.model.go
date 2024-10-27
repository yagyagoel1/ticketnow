package models

import "time"

type Show struct {
	Id           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"not null" json:"name"`
	Description  string `gorm:"not null" json:"description"`
	Image        string `gorm:"not null" json:"image"`
	TicketTypes  []TicketType
	Bookings     []Booking
	BookingLocks []BookingLock
	Location     string    `gorm:"not null" json:"location"`
	ShowTiming   time.Time `gorm:"not null" json:"showTiming"`
}

type TicketType struct {
	Id          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"not null" json:"name"`
	Price       float64 `gorm:"not null" json:"price"`
	Count       uint    `gorm:"not null" json:"count"`
	ShowId      uint    `json:"showId" gorm:"not null"`
	TicketCount []TicketCount
	Show        Show `gorm:"foreignKey:ShowId"`
}
