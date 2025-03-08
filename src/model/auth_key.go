package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthKey struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Key       string             `json:"key" bson:"key"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
