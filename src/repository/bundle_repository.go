package repository

import (
	"context"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BundleRepository interface {
	CreateBundle(ctx context.Context, bundle *model.Bundle) (*model.Bundle, error)
}

type bundleRepositoryImpl struct {
	Connection *mongo.Database
}

func NewBundleRepository(db *mongo.Database) BundleRepository {
	return &bundleRepositoryImpl{Connection: db}
}

func (bundleRepository *bundleRepositoryImpl) CreateBundle(ctx context.Context, bundle *model.Bundle) (*model.Bundle, error) {
	bundle.CreatedAt = time.Now()
	bundle.UpdatedAt = time.Now()
	collection := bundleRepository.Connection.Collection("bundles")
	_, err := collection.InsertOne(ctx, &bundle, options.InsertOne())
	if err != nil {
		return nil, err
	}
	return bundle, nil
}
