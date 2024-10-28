package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id           uint   `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Name         string `json:"name" gorm:"not null"`
	Email        string `json:"email" gorm:"not null;unique"`
	Password     string `json:"password" gorm:"not null"`
	Token        string `json:"token" gorm:"not null"`
	IsAdmin      bool   `json:"isAdmin" gorm:"not null;default:false"`
	Verified     bool   `json:"verified" gorm:"not null;default:false"`
	CreateEvent  bool   `json:"createEvent" gorm:"not null;default:false"`
	Bookings     []Booking
	Shows        []Show
	BookingLocks []BookingLock
}

func (u *User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if len(u.Password) > 0 {
		hashedPassword, err := u.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}
