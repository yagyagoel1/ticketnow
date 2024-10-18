package models

import "time"

type Booking struct {
	Id           uint       `gorm:"primaryKey;autoincrement" json:"id"`
	UserId       uint       `json:"userId" grom:"not null"`
	User         User       `gorm:"foreignKey:UserId"`
	ShowId       uint       `json:"showId" grom:"not null"`
	Show         Show       `grom:"foreignKey:ShowId"`
	TicketTypeId uint       `json:"ticketTypeId" grom:"not null"`
	TicketType   TicketType `gorm:"foreignKey:TicketType"`
	TotalCost    float64    `json:"totalCost" grom:"not null"`
}

type BookingLock struct {
	Id           uint       `gorm:"primaryKey;autoincrement" json:"id"`
	LockTime     time.Time  `gorm:"not null" json:"lockTime"`
	UserId       uint       `json:"userId" grom:"not null"`
	User         User       `gorm:"foreignKey:UserId"`
	ShowId       uint       `json:"showId" grom:"not null"`
	Show         Show       `grom:"foreignKey:ShowId"`
	TicketTypeId uint       `json:"ticketTypeId" grom:"not null"`
	TicketType   TicketType `gorm:"foreignKey:TicketType"`
	TotalCost    float64    `json:"totalCost" grom:"not null"`
}
