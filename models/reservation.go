package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Reservation predstavlja jednu rezervaciju Usera
type Reservation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username      string             `bson:"username,omitempty" json:"username,omitempty"`
	UserID        primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	MovieID       primitive.ObjectID `bson:"movieId,omitempty" json:"movieId,omitempty"`
	RepertoiresID primitive.ObjectID `bson:"repertoiresId,omitempty" json:"repertoiresId,omitempty"`
	MovieTitle    string             `bson:"movieTitle,omitempty" json:"movieTitle,omitempty"`
	Date          time.Time          `bson:"date,omitempty" json:"date,omitempty"`
	Time          string             `bson:"time,omitempty" json:"time,omitempty"`
	Hall          string             `bson:"hall,omitempty" json:"hall,omitempty"`
	CreationDate  time.Time          `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
	ReservSeats   []string           `bson:"reservSeats,omitempty" json:"reservSeats,omitempty"`
}
