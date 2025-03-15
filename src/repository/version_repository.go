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

type VersionRepository interface {
	Create(ctx context.Context, version *model.Version) (*model.Version, error)
	GetByEnvironmentIdAndAppVersion(ctx context.Context, deploymentId primitive.ObjectID, appVersion string) (*model.Version, error)
	UpdateCurrentBundleId(ctx context.Context, id primitive.ObjectID, currentBundleId primitive.ObjectID) (*model.Version, error)
	GetByEnvironmentAndVersion(ctx context.Context, environment string, version string) (*model.Version, error)
	GetLatestVersionByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) (*model.Version, error)
}

type versionRepository struct {
	Connection *mongo.Database
}

func NewVersionRepository(db *mongo.Database) VersionRepository {
	return &versionRepository{Connection: db}
}

func (v *versionRepository) Create(ctx context.Context, version *model.Version) (*model.Version, error) {
	version.CreatedAt = time.Now()
	version.UpdatedAt = time.Now()
	collection := v.Connection.Collection("versions")
	insertedVersion, err := collection.InsertOne(ctx, version)
	if err != nil {
		return nil, err
	}
	version.Id = insertedVersion.InsertedID.(primitive.ObjectID)
	return version, nil
}

func (v *versionRepository) GetByEnvironmentIdAndAppVersion(ctx context.Context, environmentId primitive.ObjectID, appVersion string) (*model.Version, error) {
	var version model.Version
	collection := v.Connection.Collection("versions")
	filter := bson.M{"environmentId": environmentId, "appVersion": appVersion}
	err := collection.FindOne(ctx, filter).Decode(&version)
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (v *versionRepository) UpdateCurrentBundleId(ctx context.Context, id primitive.ObjectID, currentBundleId primitive.ObjectID) (*model.Version, error) {
	collection := v.Connection.Collection("versions")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"currentBundleId": currentBundleId, "updatedAt": time.Now()}})
	if err != nil {
		return nil, err
	}
	return &model.Version{Id: id, CurrentBundleId: currentBundleId, UpdatedAt: time.Now()}, nil
}

func (v *versionRepository) GetByEnvironmentAndVersion(ctx context.Context, environment string, version string) (*model.Version, error) {
	collection := v.Connection.Collection("versions")
	var versionDocument model.Version
	err := collection.FindOne(ctx, bson.M{"environment": environment, "version": version}).Decode(&versionDocument)
	if err != nil {
		return nil, err
	}
	return &versionDocument, nil
}

func (v *versionRepository) GetLatestVersionByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) (*model.Version, error) {
	collection := v.Connection.Collection("versions")
	var versionDocument model.Version
	opts := options.FindOne().SetSort(bson.D{{Key: "versionNumber", Value: -1}})
	err := collection.FindOne(ctx, bson.M{"environmentId": environmentId}, opts).Decode(&versionDocument)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &versionDocument, nil
}
