package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username         string             `json:"username" bson:"username"`
	Password         string             `json:"password" bson:"password"`
	DateOfCreation   time.Time          `bson:"creation_date,omitempty" json:"creation_date,omitempty"`
	DateOfLastUpdate time.Time          `bson:"update_date,omitempty" json:"update_date,omitempty"`
	Roles            []string           `bson:"roles,omitempty" json:"roles,omitempty"`
}
