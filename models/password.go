package models

import "gorm.io/gorm"

type PasswordResetToken struct {
	gorm.Model
	UID   string
	Token string `gorm:"unique"`
}
