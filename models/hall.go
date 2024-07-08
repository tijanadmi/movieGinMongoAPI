package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Hall predstavlja podatke o bioskopskoj sali
type Hall struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	Rows      []string           `bson:"rows,omitempty" json:"rows,omitempty"`
	Cols      []int              `bson:"cols,omitempty" json:"cols,omitempty"`
	CreatedAt time.Time          `bson:"creation_date,omitempty" json:"creation_date,omitempty"`
}

// type ListLimitOffsetParams struct {
// 	Limit  int32  `json:"limit"`
// 	Offset int32  `json:"offset"`
// }
