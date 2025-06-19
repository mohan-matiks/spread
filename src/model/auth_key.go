package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthKey struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Key       string             `json:"key" bson:"key"`
	IsValid   bool               `json:"isValid" bson:"isValid,default:true"`
	CreatedBy string             `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
