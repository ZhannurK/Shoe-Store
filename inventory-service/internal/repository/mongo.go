package repository

import (
	"context"
	"fmt"

	"github.com/shoe-store/inventory-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("sneakers"),
	}
}

func (r *MongoRepository) CreateSneaker(ctx context.Context, sneaker *models.Sneaker) (*models.Sneaker, error) {
	result, err := r.collection.InsertOne(ctx, sneaker)
	if err != nil {
		return nil, fmt.Errorf("failed to insert sneaker: %w", err)
	}
	sneaker.ID = result.InsertedID.(primitive.ObjectID)
	return sneaker, nil
}

func (r *MongoRepository) GetSneakers(ctx context.Context, page, limit int32) ([]*models.Sneaker, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}
	skip := (page - 1) * limit
	skip64 := int64(skip)
	limit64 := int64(limit)

	option := &options.FindOptions{
		Skip:  &skip64,
		Limit: &limit64,
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, option)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var sneakers []*models.Sneaker
	if err = cursor.All(ctx, &sneakers); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return sneakers, total, nil
}

func (r *MongoRepository) GetPublicSneakers(ctx context.Context, page, limit int32) ([]*models.Sneaker, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}
	skip := (page - 1) * limit
	skip64 := int64(skip)
	limit64 := int64(limit)

	option := &options.FindOptions{
		Skip:  &skip64,
		Limit: &limit64,
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, option)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var sneakers []*models.Sneaker
	if err = cursor.All(ctx, &sneakers); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return sneakers, total, nil
}

func (r *MongoRepository) UpdateSneaker(ctx context.Context, id primitive.ObjectID, sneaker *models.Sneaker) error {
	update := bson.M{
		"$set": bson.M{
			"brand": sneaker.Brand,
			"model": sneaker.Model,
			"price": sneaker.Price,
			"color": sneaker.Color,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *MongoRepository) DeleteSneaker(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoRepository) GetSneakerByID(ctx context.Context, id primitive.ObjectID) (*models.Sneaker, error) {
	var sneaker models.Sneaker
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sneaker)
	if err != nil {
		return nil, err
	}
	return &sneaker, nil
}

func (r *MongoRepository) UpdateSneakerStock(ctx context.Context, id primitive.ObjectID, newStock int) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"stock": newStock}},
	)
	return err
}
