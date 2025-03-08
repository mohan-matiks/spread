package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type App struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	OS        string             `json:"os" bson:"os"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
