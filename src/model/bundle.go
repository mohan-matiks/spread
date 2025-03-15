package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bundle struct {
	Id            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EnvironmentId primitive.ObjectID `json:"environmentId" bson:"environmentId"`
	VersionId     primitive.ObjectID `json:"versionId" bson:"versionId"`
	AppId         primitive.ObjectID `json:"appId" bson:"appId"`
	Size          int64              `json:"size" bson:"size"`
	SequenceId    int64              `json:"sequenceId" bson:"sequenceId"`
	Hash          string             `json:"hash" bson:"hash"`
	DownloadFile  string             `json:"downloadFile" bson:"downloadFile"`
	IsMandatory   bool               `json:"isMandatory" bson:"isMandatory" default:"false"`
	Failed        int                `json:"failed" bson:"failed" default:"0"`
	Installed     int                `json:"installed" bson:"installed" default:"0"`
	Active        int                `json:"active" bson:"active" default:"0"`
	Description   string             `json:"description" bson:"description"`
	Label         string             `json:"label" bson:"label"`
	IsValid       bool               `json:"isValid" bson:"isValid" default:"true"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}
