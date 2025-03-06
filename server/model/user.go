package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FullName     string             `bson:"_fullName,omitempty" json:"_fullName,omitempty"`
	Email        string             `bson:"_email,omitempty" json:"_email,omitempty"`
	HashPassword string             `bson:"_password,omitempty" json:"_password,omitempty"`
}
