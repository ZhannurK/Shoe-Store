package repositories

import (
	"auth-service/internal/entities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByToken(ctx context.Context, token string) (*entities.User, error)
	GetUserByID(id string) (*entities.User, error)
}

var _ UserRepository = (*UserMongoRepo)(nil)

type UserMongoRepo struct {
	coll *mongo.Collection
}

func NewUserMongoRepo(db *mongo.Client, dbName, collName string) *UserMongoRepo {
	return &UserMongoRepo{
		coll: db.Database(dbName).Collection(collName),
	}
}

func (r *UserMongoRepo) GetUserByID(id string) (*entities.User, error) {
	var user entities.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	err = r.coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserMongoRepo) Create(ctx context.Context, user *entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(ctx, user)
	return err
}

func (r *UserMongoRepo) Update(ctx context.Context, user *entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	_, err := r.coll.UpdateOne(ctx, filter, bson.M{"$set": user})
	return err
}

func (r *UserMongoRepo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
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

func (r *UserMongoRepo) FindByToken(ctx context.Context, token string) (*entities.User, error) {
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
func (r *UserMongoRepo) FindByID(ctx context.Context, id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user entities.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
