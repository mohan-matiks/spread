package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Version struct {
	Id              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EnvironmentId   primitive.ObjectID `json:"environmentId" bson:"environmentId"`
	AppVersion      string             `json:"appVersion" bson:"appVersion"`
	VersionNumber   int64              `json:"versionNumber" bson:"versionNumber"`
	CurrentBundleId primitive.ObjectID `json:"currentBundleId" bson:"currentBundleId"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
}
