package users

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"firstName" bson:"firstName"`
	Surname    string             `json:"secondName" bson:"secondName"`
	Patronymic string             `json:"patronymic" bson:"patronymic"`
	Type       string             `json:"type" bson:"type"`
	Branch     string             `json:"branch" bson:"branch"`
	Department string             `json:"department" bson:"department"`
	Position   string             `json:"position" bson:"position"`
}

func (u *User) ToBson() bson.M {
	return bson.M{
		"name":       u.Name,
		"surname":    u.Surname,
		"patronymic": u.Patronymic,
		"type":       u.Type,
		"branch":     u.Branch,
		"department": u.Department,
		"position":   u.Position,
	}
}
