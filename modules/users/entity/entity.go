package entity

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	//Id               primitive.ObjectID `json:"id" bson:"_id"`
	Surname    string             `json:"surname" bson:"surname"`
	Name       string             `json:"name" bson:"name"`
	Patronymic string             `json:"patronymic" bson:"patronymic"`
	Branch     primitive.ObjectID `json:"branch" bson:"branch"`
	Department primitive.ObjectID `json:"department" bson:"department"`
	Position   primitive.ObjectID `json:"position" bson:"position"`
	StartDate  string             `json:"start_date" bson:"start_date"`
	IsAdmin    bool               `json:"is_admin" bson:"is_admin"`
}

type UserFilters struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Surname    string             `json:"surname,omitempty" bson:"surname,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Patronymic string             `json:"patronymic,omitempty" bson:"patronymic,omitempty"`
	Branch     primitive.ObjectID `json:"branch,omitempty" bson:"branch,omitempty"`
	Department primitive.ObjectID `json:"department,omitempty" bson:"department,omitempty"`
	Position   primitive.ObjectID `json:"position,omitempty" bson:"position,omitempty"`
	StartDate  string             `json:"start_date,omitempty" bson:"start_date,omitempty"`
	IsAdmin    bool               `json:"is_admin,omitempty" bson:"is_admin,omitempty"`
}
