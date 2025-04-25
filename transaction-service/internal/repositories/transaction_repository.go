package repositories

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"transaction-service/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	UpdateStatus(ctx context.Context, transactionID string, status domain.TransactionStatus) error
	GetByID(ctx context.Context, transactionID string) (*domain.Transaction, error)
	Delete(ctx context.Context, transactionID string) error
}

type transactionRepo struct {
	collection *mongo.Collection
}

//func NewTransactionRepository(client *mongo.Client, dbName string) TransactionRepository {
//	return &transactionRepo{
//		collection: client.Database(dbName).Collection("transactions"),
//	}
//}

func NewTransactionRepository(client *mongo.Client, dbName string) TransactionRepository {
	return &transactionRepo{
		collection: client.Database(dbName).Collection("transactions"),
	}
}

//func (r *transactionRepo) Create(ctx context.Context, tx *domain.Transaction) error {
//	tx.CreatedAt = time.Now()
//	tx.UpdatedAt = time.Now()
//	_, err := r.collection.InsertOne(ctx, tx)
//	return err
//}

func (r *transactionRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	now := time.Now()

	if tx.ID.IsZero() {
		tx.ID = primitive.NewObjectID()
	}
	if tx.Status == "" {
		tx.Status = domain.StatusPending
	}
	tx.CreatedAt = now
	tx.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, tx)
	return err
}

//func (r *transactionRepo) UpdateStatus(ctx context.Context, transactionID string, status domain.TransactionStatus) error {
//	filter := bson.M{"transactionId": transactionID}
//	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}}
//	_, err := r.collection.UpdateOne(ctx, filter, update)
//	return err
//}

func (r *transactionRepo) UpdateStatus(
	ctx context.Context,
	transactionID string,
	status domain.TransactionStatus,
) error {

	filter := bson.M{"transactionId": transactionID}
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

//func (r *transactionRepo) GetByID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
//	var tx domain.Transaction
//	err := r.collection.FindOne(ctx, bson.M{"transactionId": transactionID}).Decode(&tx)
//	if err != nil {
//		return nil, err
//	}
//	return &tx, nil
//}

func (r *transactionRepo) GetByID(
	ctx context.Context,
	transactionID string,
) (*domain.Transaction, error) {

	filter := bson.M{"transactionId": transactionID}
	var tx domain.Transaction

	err := r.collection.FindOne(ctx, filter).Decode(&tx)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %w", err)
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
