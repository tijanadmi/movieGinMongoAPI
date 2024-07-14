package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrRecordNotFound = mongo.ErrNoDocuments
