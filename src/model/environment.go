package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Environment struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AppId     primitive.ObjectID `json:"appId" bson:"appId"`
	Name      string             `json:"name" bson:"name"`
	Key       string             `json:"key" bson:"key"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
