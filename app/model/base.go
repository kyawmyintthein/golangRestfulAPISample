package model

import "github.com/mongodb/mongo-go-driver/bson/objectid"

type MongoBaseModel struct{
	RawId             *objectid.ObjectID `json:"raw_id,omitempty" bson:"_id,omitempty"`
	ID                string             `json:"id,omitempty" bson:"id,omitempty"`
	CreatedAt         int64              `json:"created_at" bson:"created_at"`
	UpdatedAt         int64              `json:"updated_at" bson:"updated_at"`
}