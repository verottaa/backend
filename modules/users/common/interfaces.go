package common

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"verottaa/modules/users/entity"
)

type Reader interface {
	Find(id primitive.ObjectID) (*entity.User, error)
	FindAll() ([]entity.User, error)
}

type Writer interface {
	Update(id primitive.ObjectID, user *entity.User) error
	Store(user *entity.User) (primitive.ObjectID, error)
	Delete(id primitive.ObjectID) error
	DeleteMany(filter entity.UserFilters) (int64, error)
	DeleteAll() (int64, error)
}
