package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}
