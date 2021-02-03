package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"verottaa/modules/plans/entity"
)

var once sync.Once
var destroyCh = make(chan bool)
var getCollectionNameCh = make(chan chan string)

var repositoryInstance *repository

type repository struct {
	collectionName string
}

func GetRepository() Repository {
	once.Do(func() {
		repositoryInstance = createRepository()
		go func() {
			for
			{
				select {
				case ch := <-getCollectionNameCh:
					ch <- repositoryInstance.collectionName
				case <-destroyCh:
					return
				}
			}
		}()
	})

	return repositoryInstance
}

func createRepository() *repository {
	repo := new(repository)
	repo.collectionName = "users"
	return repo
}

func (r repository) Destroy() {
	destroyCh <- true
	close(destroyCh)
	repositoryInstance = nil
}

func (r repository) GetCollectionName() string {
	resCh := make(chan string)
	defer close(resCh)
	getCollectionNameCh <- resCh
	return <-resCh
}

func (r repository) Find(id primitive.ObjectID) (*entity.Plan, error) {
	panic("Not implemented")
}

func (r repository) FindAll() ([]*entity.Plan, error) {
	panic("Not implemented")
}

func (r repository) Update(id primitive.ObjectID, plan *entity.Plan) error {
	panic("Not implemented")
}

func (r repository) Store(plan *entity.Plan) (primitive.ObjectID, error) {
	panic("Not implemented")
}

func (r repository) Delete(id primitive.ObjectID) error {
	panic("Not implemented")
}

func (r repository) DeleteMany(filter entity.PlanFilters) (int64, error) {
	panic("Not implemented")
}

func (r repository) DeleteAll() (int64, error) {
	panic("Not implemented")
}