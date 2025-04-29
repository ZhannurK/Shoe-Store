package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sneaker struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Brand string             `json:"brand" bson:"brand"`
	Model string             `json:"model" bson:"model"`
	Price int                `json:"price" bson:"price"`
	Color string             `json:"color" bson:"color"`
}
