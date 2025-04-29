package repositories

import (
	"auth-service/internal/entities"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByToken(ctx context.Context, token string) (*entities.User, error)
}

type userMongoRepo struct {
	coll *mongo.Collection
}

func NewUserMongoRepo(db *mongo.Client, dbName, collName string) UserRepository {
	return &userMongoRepo{
		coll: db.Database(dbName).Collection(collName),
	}
}

func (r *userMongoRepo) Create(ctx context.Context, user *entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(ctx, user)
	return err
}

func (r *userMongoRepo) Update(ctx context.Context, user *entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	_, err := r.coll.UpdateOne(ctx, filter, bson.M{"$set": user})
	return err
}

func (r *userMongoRepo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user entities.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userMongoRepo) FindByToken(ctx context.Context, token string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user entities.User
	err := r.coll.FindOne(ctx, bson.M{"confirmationToken": token}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}
	return &user, nil
}
