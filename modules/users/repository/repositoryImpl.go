package repository

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"verottaa/modules/users/entity"
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

func (r repository) Find(id primitive.ObjectID) (*entity.User, error) {
	user := &entity.User{}
	coll := mgm.Coll(user)
	err := coll.FindByID(id, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r repository) FindAll() ([]entity.User, error) {
	user := &entity.User{}
	var users []entity.User
	coll := mgm.Coll(user)
	err := coll.SimpleFind(&users, bson.M{})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r repository) Update(id primitive.ObjectID, user *entity.User) error {
	// TODO: найти способ починить ошибку апдейта
	filters := entity.UserFilters{Id: id}
	_, err := mgm.Coll(user).UpdateOne(mgm.Ctx(), filters, user)
	return err
}

func (r repository) Store(user *entity.User) (primitive.ObjectID, error) {
	err := mgm.Coll(user).Create(user)
	if err != nil {
		return user.ID, err
	}
	return user.ID, nil
}

func (r repository) Delete(id primitive.ObjectID) error {
	user := &entity.User{}
	user.SetID(id)
	return mgm.Coll(user).Delete(user)
}

func (r repository) DeleteMany(filter entity.UserFilters) (int64, error) {
	panic("Not implemented")
}

func (r repository) DeleteAll() (int64, error) {
	user := &entity.User{}
	context := mgm.Ctx()
	deleteResult, err := mgm.Coll(user).DeleteMany(context, bson.M{})
	if err != nil {
		return -1, err
	}
	return deleteResult.DeletedCount, nil
}
