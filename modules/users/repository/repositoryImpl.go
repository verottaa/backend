package repository

import (
	"errors"
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

func decodeBson(bsonObject interface{}, target interface{}) {
	m, _ := bson.Marshal(bsonObject)
	_ = bson.Unmarshal(m, target)
}

func decodeUserFromBson(bsonUser interface{}) (*entity.User, error) {
	var user entity.User
	err := user.GetUserFromBson(bsonUser)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r repository) Find(id primitive.ObjectID) (*entity.User, error) {
	filter := entity.UserFilters{
		Id: id,
	}
	code, bsonUser := databaseModule.FindOne(r.GetCollectionName(), filter)
	switch code {
	case databaseModule.FOUND:
		return decodeUserFromBson(bsonUser)
	case databaseModule.ERROR:
		return nil, errors.New("unexpected error with database connection")
	default:
		return nil, errors.New("unexpected code")
	}
}

func (r repository) FindAll() ([]*entity.User, error) {
	filter := entity.UserFilters{}
	code, bsonUsers := databaseModule.FindMany(r.GetCollectionName(), filter)
	switch code {
	case databaseModule.FOUND_ANY:
		var users []*entity.User
		for _, bsonUser := range bsonUsers {
			user, err := decodeUserFromBson(bsonUser)
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		return users, nil
	case databaseModule.ERROR:
		return nil, errors.New("unexpected error with database connection")
	default:
		return nil, errors.New("unexpected code")
	}
}

func (r repository) Update(user *entity.User) error {
	filter:= entity.UserFilters{
		Id: user.Id,
	}
	code := databaseModule.UpdateOne(r.GetCollectionName(), filter, user.ToUpdateObjectData())
	switch code {
	case databaseModule.UPDATED:
		return nil
	case databaseModule.NOT_FOUND:
		return errors.New("object didn't found in database")
	case databaseModule.ERROR:
		return errors.New("unexpected error with database connection")
	default:
		return errors.New("unexpected code")
	}
}

func (r repository) Store(user *entity.User) (primitive.ObjectID, error) {
	user.Id = databaseModule.GenerateObjectID()
	code, bsonId := databaseModule.PushOne(r.GetCollectionName(), user)
	switch code {
	case databaseModule.CREATED:
		var id primitive.ObjectID
		decodeBson(bsonId, &id)

		if user.Id.String() == id.String() {
			return id, nil
		}
		return primitive.ObjectID{}, errors.New("validation didn't pass")
	case databaseModule.ERROR:
		return primitive.ObjectID{}, errors.New("unexpected error with database connection")
	default:
		return primitive.ObjectID{}, errors.New("unexpected code")

	}
}

func (r repository) Delete(id primitive.ObjectID) error {
	filter := entity.UserFilters{
		Id: id,
	}
	code := databaseModule.DeleteOne(r.GetCollectionName(), filter)
	switch code {
	case databaseModule.DELETED:
		return nil
	case databaseModule.NOT_FOUND:
		return errors.New("object didn't found in database")
	case databaseModule.ERROR:
		return errors.New("unexpected error with database connection")
	default:
		return errors.New("unexpected code")
	}
}

func (r repository) DeleteMany(filter entity.UserFilters) (int64, error) {
	code, quantity := databaseModule.DeleteMany(r.GetCollectionName(), filter)
	switch code {
	case databaseModule.DELETED:
		return quantity, nil
	case databaseModule.NOT_FOUND:
		return quantity, errors.New("object didn't found in database")
	case databaseModule.ERROR:
		return quantity, errors.New("unexpected error with database connection")
	default:
		return quantity, errors.New("unexpected code")
	}
}

func (r repository) DeleteAll() (int64, error) {
	filter := entity.UserFilters{}
	return r.DeleteMany(filter)
}
