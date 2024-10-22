package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id          uint      `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Name        string    `json:"name" gorm:"not null"`
	Email       string    `json:"email" gorm:"not null"`
	Password    string    `json:"password" gorm:"not null"`
	Token       string    `json:"token" gorm:"not null"`
	Role        string    `json:"role" gorm:"not null"`
	Verified    time.Time `json:"verified" gorm:"not null"`
	CreateEvent time.Time `json:"CreateEvent" gorm:"not null"`
	Bookings    []Booking
	Shows       []Show
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
