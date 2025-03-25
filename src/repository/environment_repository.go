package repository

import (
	"context"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EnvironmentRepository interface {
	Insert(ctx context.Context, environment *model.Environment) (*model.Environment, error)
	GetByKey(ctx context.Context, key string) (*model.Environment, error)
	GetByAppIdAndName(ctx context.Context, appId primitive.ObjectID, name string) (*model.Environment, error)
	GetAllByAppId(ctx context.Context, appId primitive.ObjectID) ([]*model.Environment, error)
	GetByIdAndAppId(ctx context.Context, id primitive.ObjectID, appId primitive.ObjectID) (*model.Environment, error)
}

type environmentRepositoryImpl struct {
	Connection *mongo.Database
}

func NewEnvironmentRepository(db *mongo.Database) EnvironmentRepository {
	return &environmentRepositoryImpl{Connection: db}
}

func (environmentRepository *environmentRepositoryImpl) Insert(ctx context.Context, environment *model.Environment) (*model.Environment, error) {
	environment.CreatedAt = time.Now()
	environment.UpdatedAt = time.Now()
	collection := environmentRepository.Connection.Collection("environments")
	_, err := collection.InsertOne(ctx, &environment, options.InsertOne())
	if err != nil {
		return nil, err
	}
	return environment, nil
}

func (environmentRepository *environmentRepositoryImpl) GetByKey(ctx context.Context, key string) (*model.Environment, error) {
	var environment model.Environment
	collection := environmentRepository.Connection.Collection("environments")
	filter := bson.M{"key": key}
	err := collection.FindOne(ctx, filter).Decode(&environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}

func (environmentRepository *environmentRepositoryImpl) GetByAppIdAndName(ctx context.Context, appId primitive.ObjectID, name string) (*model.Environment, error) {
	var environment model.Environment
	collection := environmentRepository.Connection.Collection("environments")
	filter := bson.M{"appId": appId, "name": name}
	err := collection.FindOne(ctx, filter).Decode(&environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}

func (environmentRepository *environmentRepositoryImpl) GetAllByAppId(ctx context.Context, appId primitive.ObjectID) ([]*model.Environment, error) {
	collection := environmentRepository.Connection.Collection("environments")
	filter := bson.M{"appId": appId}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var environments []*model.Environment
	for cursor.Next(ctx) {
		var environment model.Environment
		if err := cursor.Decode(&environment); err != nil {
			return nil, err
		}
		environments = append(environments, &environment)
	}
	return environments, nil
}

func (environmentRepository *environmentRepositoryImpl) GetByIdAndAppId(ctx context.Context, id primitive.ObjectID, appId primitive.ObjectID) (*model.Environment, error) {
	collection := environmentRepository.Connection.Collection("environments")
	filter := bson.M{"_id": id, "appId": appId}
	var environment model.Environment
	err := collection.FindOne(ctx, filter).Decode(&environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}
