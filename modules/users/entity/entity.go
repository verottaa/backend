package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Surname    string             `json:"surname" bson:"surname"`
	Patronymic string             `json:"patronymic" bson:"patronymic"`
	Type       string             `json:"type" bson:"type"`
	Branch     string             `json:"branch" bson:"branch"`
	Department string             `json:"department" bson:"department"`
	Position   string             `json:"position" bson:"position"`
	StartDate  string             `json:"start_date" bson:"start_date"`
	IsAdmin    bool               `json:"is_admin" bson:"is_admin"`
}

type UserFilters struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Surname    string             `json:"surname,omitempty" bson:"surname,omitempty"`
	Patronymic string             `json:"patronymic,omitempty" bson:"patronymic,omitempty"`
	Type       string             `json:"type,omitempty" bson:"type,omitempty"`
	Branch     string             `json:"branch,omitempty" bson:"branch,omitempty"`
	Department string             `json:"department,omitempty" bson:"department,omitempty"`
	Position   string             `json:"position,omitempty" bson:"position,omitempty"`
	StartDate  string             `json:"start_date,omitempty" bson:"start_date,omitempty"`
	IsAdmin    bool               `json:"is_admin,omitempty" bson:"is_admin,omitempty"`
}
