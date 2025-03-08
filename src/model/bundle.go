package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bundle struct {
	Id                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeploymentId        primitive.ObjectID `json:"deploymentId" bson:"deploymentId"`
	DeploymentVersionId primitive.ObjectID `json:"deploymentVersionId" bson:"deploymentVersionId"`
	BundleVersionId     int64              `json:"bundleVersionId" bson:"bundleVersionId"`
	Size                int64              `json:"size" bson:"size"`
	Hash                string             `json:"hash" bson:"hash"`
	Download            string             `json:"download" bson:"download"`
	IsMandatory         bool               `json:"isMandatory" bson:"isMandatory"`
	Active              int                `json:"active" bson:"active" default:"0"`
	Failed              int                `json:"failed" bson:"failed" default:"0"`
	Installed           int                `json:"installed" bson:"installed" default:"0"`
	Description         string             `json:"description" bson:"description"`
	IsValid             bool               `json:"isValid" bson:"isValid"`
	CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt" bson:"updatedAt"`
}
