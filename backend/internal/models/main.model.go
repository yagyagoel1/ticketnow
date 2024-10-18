package models

import "gorm.io/gorm"

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Show{}, &Booking{}, &TicketType{}, &BookingLock{})
	return err
}
