package repository

import (
	"context"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Insert(ctx context.Context, user *model.User) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*model.User, error)
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Insert(ctx context.Context, user *model.User) (*model.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	collection := r.db.Collection("users")
	createdUser, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = createdUser.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	collection := r.db.Collection("users")
	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetById(ctx context.Context, id primitive.ObjectID) (*model.User, error) {

	collection := r.db.Collection("users")
	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
