package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	UID             primitive.ObjectID `bson:"_id,omitempty"`
	Email           string             `bson:"email" json:"email"`
	Password        string             `bson:"password" json:"password"`
	CustomerName    string             `bson:"customer_name" json:"customer_name"`
	MobileNumber    string             `bson:"mobile_number" json:"mobile_number"`
	DeliveryAddress string             `bson:"delivery_address" json:"delivery_address"`
}

type Merchant struct {
	UID          primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"password"`
	MerchantName string             `bson:"merchant_name" json:"merchant_name"`
	MobileNumber string             `bson:"mobile_number" json:"mobile_number"`
	StoreAddress string             `bson:"store_address" json:"store_address"`
}

// type DeliveryAgent struct {
// 	UID          primitive.ObjectID `bson:"_id,omitempty"`
// 	Email        string             `bson:"email" json:"email"`
// 	Password     string             `bson:"password" json:"password"`
// 	AgentName    string             `bson:"agent_name" json:"agent_name"`
// 	MobileNumber string             `bson:"mobile_number" json:"mobile_number"`
// 	Verified     bool               `bson:"verified" json:"verified"`
// }
