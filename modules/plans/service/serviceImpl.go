package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"verottaa/modules/plans/entity"
	"verottaa/modules/plans/repository"
)

var once sync.Once
var destroyCh = make(chan bool)

var serviceInstance *service

type service struct {
	repo repository.Repository
}

func GetService() Service {
	once.Do(func() {
		serviceInstance = createService()
		go func() {
			for
			{
				select {
				case <-destroyCh:
					return
				}
			}
		}()
	})

	return serviceInstance
}

func createService() *service {
	var serviceInst = new(service)
	serviceInst.repo = repository.GetRepository()
	return serviceInst
}

func (s service) Destroy() {
	destroyCh <- true
	close(destroyCh)
	serviceInstance = nil
}

func (s service) Find(id primitive.ObjectID) (*entity.Plan, error) {
	return s.repo.Find(id)
}

func (s service) FindAll() ([]*entity.Plan, error) {
	return s.repo.FindAll()
}

func (s service) Update(user *entity.Plan) error {
	return s.repo.Update(user)
}

func (s service) Store(user *entity.Plan) (primitive.ObjectID, error) {
	return s.repo.Store(user)
}

func (s service) Delete(id primitive.ObjectID) error {
	return s.repo.Delete(id)
}

func (s service) DeleteMany(filter entity.PlanFilters) (int64, error) {
	return s.repo.DeleteMany(filter)
}

func (s service) DeleteAll() (int64, error) {
	return s.repo.DeleteAll()
}
