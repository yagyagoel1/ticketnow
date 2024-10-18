package models

import "time"

type Show struct {
	Id          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
	Image       string `grom:"not null" json:"image"`
	TicketTypes []TicketType
	Bookings    []Booking
	Location    string    `gorm:"not null" json:"location"`
	ShowTiming  time.Time `grom:"not null" json:"showTiming"`
}

type TicketType struct {
	Id       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string  `gorm:"not null" json:"name"`
	Price    float64 `gorm:"not null" json:"price"`
	Count    uint    `gorm:"not null" json:"count"`
	ShowId   uint    `json:"showid" grom:"not null"`
	Bookings []Booking
	Show     Show `gorm:"foreignKey:ShowId"`
}
