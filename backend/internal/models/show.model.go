package models

import "time"

type Show struct {
	Id           uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string       `gorm:"not null" json:"name"`
	Description  string       `gorm:"not null" json:"description"`
	Image        string       `gorm:"not null" json:"image"`
	TicketTypes  []TicketType `json:"ticketTypes" gorm:"foreignKey:ShowId;constraint:OnDelete:CASCADE;"`
	Bookings     []Booking
	BookingLocks []BookingLock
	User         User      `json:"-" gorm:"foreignKey:UserId;"`
	UserId       uint      `json:"userId" gorm:"not null"`
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
	Show        Show `json:"-" gorm:"foreignKey:ShowId;constraint:OnDelete:CASCADE;"`
}
