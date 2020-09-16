package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `bson:"_id" ,json:"id"`
	FirstName  string             `bson:"firstName" ,json:"firstName"`
	SecondName string             `bson:"secondName" ,json:"secondName"`
	Patronymic string             `bson:"patronymic" ,json:"patronymic"`
	//Login      string             `bson:"login" ,json:"login"`
	//Password   string             `bson:"password" ,json:"password"`
	Type       string `bson:"type" ,json:"type"`
	Branch     string `bson:"branch" ,json:"branch"`
	Department string `bson:"department" ,json:"department"`
}
