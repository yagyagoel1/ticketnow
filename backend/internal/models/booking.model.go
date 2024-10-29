package models

import "time"

type Booking struct {
	Id          uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      uint `json:"userId" gorm:"not null"`
	User        User `gorm:"foreignKey:UserId"`
	ShowId      uint `json:"showId" gorm:"not null"`
	Show        Show `gorm:"foreignKey:ShowId"`
	TicketCount []TicketCount
	TotalCost   float64 `json:"totalCost" gorm:"not null"`
}
type TicketCount struct {
	Id                  uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingId           *uint       `json:"bookingId"`
	Booking             *Booking    `gorm:"foreignKey:BookingId"`
	BookingLockId       *uint       `json:"bookingLockId" `
	BookingLock         BookingLock `gorm:"foreignKey:BookingLockId" json:"-"`
	TicketTypeId        uint        `json:"ticketTypeId" gorm:"not null"`
	TicketType          TicketType  `gorm:"foreignKey:TicketTypeId"`
	TicketCountCategory int         `json:"ticketCountCategory" gorm:"not null"`
}

type BookingLock struct {
	Id          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	LockTime    time.Time     `gorm:"not null" json:"lockTime"`
	UserId      uint          `json:"userId" gorm:"not null"`
	User        User          `gorm:"foreignKey:UserId" json:"-"`
	ShowId      uint          `json:"showId" gorm:"not null"`
	Show        Show          `gorm:"foreignKey:ShowId"`
	TicketCount []TicketCount `gorm:"foreignKey:BookingLockId" json:"ticketCount"`
	TotalCost   float64       `json:"totalCost" gorm:"not null"`
}
