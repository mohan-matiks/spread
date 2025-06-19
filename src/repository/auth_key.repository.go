package repository

import (
	"context"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthKeyRepository interface {
	Insert(authKey *model.AuthKey) (*model.AuthKey, error)
	GetById(key string) (*model.AuthKey, error)
	GetAll(ctx context.Context) ([]*model.AuthKey, error)
}

type authKeyRepository struct {
	Connection *mongo.Database
}

func NewAuthKeyRepository(db *mongo.Database) AuthKeyRepository {
	return &authKeyRepository{Connection: db}
}

func (r *authKeyRepository) Insert(authKey *model.AuthKey) (*model.AuthKey, error) {
	authKey.CreatedAt = time.Now()
	authKey.UpdatedAt = time.Now()
	collection := r.Connection.Collection("auth_keys")
	createdAuthKey, err := collection.InsertOne(context.Background(), authKey)
	if err != nil {
		return nil, err
	}
	authKey.Id = createdAuthKey.InsertedID.(primitive.ObjectID)
	return authKey, nil
}

func (r *authKeyRepository) GetById(key string) (*model.AuthKey, error) {
	collection := r.Connection.Collection("auth_keys")
	var authKey model.AuthKey
	err := collection.FindOne(context.Background(), bson.M{"key": key}).Decode(&authKey)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &authKey, nil
}

func (r *authKeyRepository) GetAll(ctx context.Context) ([]*model.AuthKey, error) {
	collection := r.Connection.Collection("auth_keys")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var authKeys []*model.AuthKey
	if err = cursor.All(ctx, &authKeys); err != nil {
		return nil, err
	}
	return authKeys, nil
}
