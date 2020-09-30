package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	SecondName string             `json:"secondName" bson:"secondName"`
	Patronymic string             `json:"patronymic" bson:"patronymic"`
	Type       string             `json:"type" bson:"type"`
	Branch     string             `json:"branch" bson:"branch"`
	Department string             `json:"department" bson:"department"`
}
