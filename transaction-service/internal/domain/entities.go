package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	pb "transaction-service/proto"
)

type TransactionStatus string

const (
	StatusPending  TransactionStatus = "Pending Payment"
	StatusPaid     TransactionStatus = "Paid"
	StatusDeclined TransactionStatus = "Declined"
)

type CartItem struct {
	SneakerID primitive.ObjectID `json:"sneakerId" bson:"sneakerId"`
	Quantity  int                `json:"quantity" bson:"quantity"`
}

type Transaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionID string             `bson:"transactionId" json:"transactionId"`
	UserID        string             `json:"userId" bson:"userId"`
	CartItems     []CartItem         `json:"cartItems" bson:"cartItems"`
	TotalAmount   float64            `json:"totalAmount" bson:"totalAmount"`
	Status        TransactionStatus  `json:"status" bson:"status"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (t *Transaction) ToProto() *pb.Transaction {
	items := make([]*pb.CartItem, len(t.CartItems))
	for i, c := range t.CartItems {
		items[i] = &pb.CartItem{
			SneakerId: c.SneakerID.Hex(),
			Quantity:  int32(c.Quantity),
		}
	}

	return &pb.Transaction{
		Id:            t.ID.Hex(),
		TransactionId: t.TransactionID,
		UserId:        t.UserID,
		CartItems:     items,
		TotalAmount:   t.TotalAmount,
		Status:        string(t.Status),
		CreatedAt:     t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     t.UpdatedAt.Format(time.RFC3339),
	}
}
