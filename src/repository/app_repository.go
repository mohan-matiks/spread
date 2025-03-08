package repository

import (
	"context"

	"github.com/SwishHQ/spread/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppRepository interface {
	InsertApp(ctx context.Context, app *model.App) (*model.App, error)
	GetAppByName(ctx context.Context, name string) (*model.App, error)
}

type appRepositoryImpl struct {
	db *mongo.Database
}

func NewAppRepository(db *mongo.Database) AppRepository {
	return &appRepositoryImpl{db: db}
}

func (appRepository *appRepositoryImpl) InsertApp(ctx context.Context, app *model.App) (*model.App, error) {
	collection := appRepository.db.Collection("apps")
	_, err := collection.InsertOne(ctx, app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (appRepository *appRepositoryImpl) GetAppByName(ctx context.Context, name string) (*model.App, error) {
	collection := appRepository.db.Collection("apps")
	var app model.App
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
