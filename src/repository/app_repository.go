package repository

import (
	"context"

	"github.com/SwishHQ/spread/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppRepository interface {
	Insert(ctx context.Context, app *model.App) (*model.App, error)
	GetByName(ctx context.Context, name string) (*model.App, error)
	GetAll(ctx context.Context) ([]*model.App, error)
}

type appRepositoryImpl struct {
	db *mongo.Database
}

func NewAppRepository(db *mongo.Database) AppRepository {
	return &appRepositoryImpl{db: db}
}

func (appRepository *appRepositoryImpl) Insert(ctx context.Context, app *model.App) (*model.App, error) {
	collection := appRepository.db.Collection("apps")
	_, err := collection.InsertOne(ctx, app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (appRepository *appRepositoryImpl) GetByName(ctx context.Context, name string) (*model.App, error) {
	collection := appRepository.db.Collection("apps")
	var app model.App
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (appRepository *appRepositoryImpl) GetAll(ctx context.Context) ([]*model.App, error) {
	collection := appRepository.db.Collection("apps")
	var apps []*model.App
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &apps); err != nil {
		return nil, err
	}
	return apps, nil
}
