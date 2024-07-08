package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repertoire predstavlja jednu projekciju filma
type Repertoire struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MovieID         primitive.ObjectID `bson:"movieId,omitempty" json:"movieId,omitempty"`
	DateSt          string             `bson:"dateSt,omitempty" json:"dateSt,omitempty"`
	Date            time.Time          `bson:"date,omitempty" json:"date,omitempty"`
	Time            string             `bson:"time,omitempty" json:"time,omitempty"`
	Hall            string             `bson:"hall,omitempty" json:"hall,omitempty"`
	NumOfTickets    int                `bson:"numOfTickets,omitempty" json:"numOfTickets,omitempty"`
	NumOfResTickets int                `bson:"numOfResTickets" json:"numOfResTickets"`
	CreatedAt       time.Time          `bson:"creation_date,omitempty" json:"creation_date,omitempty"`
	ReservSeats     []string           `bson:"reservSeats,omitempty" json:"reservSeats,omitempty"`
}
