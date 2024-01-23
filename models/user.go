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
	UID             string `gorm:"unique"`
	CustomerName    string
	MobileNumber    string
	DeliveryAddress string
}

type Merchant struct {
	gorm.Model
	UID          string `gorm:"unique"`
	MerchantName string
	MobileNumber string
	StoreAddress string
}

type DeliveryAgent struct {
	gorm.Model
	UID           string `gorm:"unique"`
	LicenseNumber string
	VehicleType   string
	VehicleNumber string
}
