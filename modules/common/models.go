package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectCreatedDto struct {
	Id primitive.ObjectID `json:"id"`
}
