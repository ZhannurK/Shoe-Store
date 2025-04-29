package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Email             string             `json:"email" bson:"email"`
	Name              string             `json:"name" bson:"name"`
	Password          string             `json:"password" bson:"password"`
	Verified          bool               `json:"verified" bson:"verified"`
	ConfirmationToken string             `json:"confirmationToken" bson:"confirmationToken"`
	Role              string             `json:"role" bson:"role"`
}
