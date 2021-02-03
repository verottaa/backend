package common

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"verottaa/modules/plans/entity"
)

// TODO: Внести в интерфейсы работу с шагами плана

type Reader interface {
	Find(id primitive.ObjectID) (*entity.Plan, error)
	FindAll() ([]*entity.Plan, error)
}

type Writer interface {
	Update(id primitive.ObjectID, user *entity.Plan) error
	Store(user *entity.Plan) (primitive.ObjectID, error)
	Delete(id primitive.ObjectID) error
	DeleteMany(filter entity.PlanFilters) (int64, error)
	DeleteAll() (int64, error)
}
