package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UID      string `gorm:"unique"`
	Email    string `gorm:"unique"`
	UserType string
	Password string
}

type Customer struct {
	gorm.Model
	UID     string `gorm:"unique"`
	Address string
}

type Merchant struct {
	gorm.Model
	UID          string `gorm:"unique"`
	MerchantName string
	StoreAddress string
}

type DeliveryAgent struct {
	gorm.Model
	UID           string `gorm:"unique"`
	LicenseNumber string
	VehicleType   string
	VehicleNumber string
}
