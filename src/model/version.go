package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Version struct {
	Id            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeploymentId  primitive.ObjectID `json:"deploymentId" bson:"deploymentId"`
	BundleId      string             `json:"bundleId" bson:"bundleId"`
	AppVersion    string             `json:"appVersion" bson:"appVersion"`
	VersionNumber int64              `json:"versionNumber" bson:"versionNumber"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
}
