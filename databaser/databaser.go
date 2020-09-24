package databaser

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"verottaa/config"
	"verottaa/constants"
	"verottaa/databaser/users"
	"verottaa/models"
)

var configuration = config.GetConfiguration()

type databaser struct {
	client          *mongo.Client
	userCollection_ users.UserCollection
}

type DB interface {
	models.Destroyable
	UserCollection
}

type UserCollection interface {
	CreateUser(user models.User) (interface{}, error)
	ReadAllUsers() ([]models.User, error)
	ReadUserById(id primitive.ObjectID) (models.User, error)
	UpdateUser(id primitive.ObjectID, user models.User) error
	DeleteUserById(id primitive.ObjectID) error
	DeleteAllUsers() error
}

var destroyCh = make(chan bool)

var instance *databaser
var once sync.Once

func initDatabaser() *databaser {
	db := new(databaser)

	var err error
	db.client, err = mongo.NewClient(options.Client().ApplyURI(configuration.GetDatabaseHost()))
	if err != nil {
		fmt.Println("[ERROR]: ", err)
	}

	db.userCollection_ = users.GetUserCollection(database)

	return db
}

func database(ctx context.Context) *mongo.Database {
	err := instance.client.Connect(ctx)
	if err != nil {
		fmt.Println("[ERROR]: ", err)
	}

	err = instance.client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return instance.client.Database(constants.DATABASE_NAME)
}

func GetDatabaser() DB {
	once.Do(func() {
		instance = initDatabaser()

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

	return instance
}

func (d databaser) Destroy() {
	destroyCh <- true
	close(destroyCh)
	instance.userCollection().Destroy()
	instance = nil
}

func (d databaser) userCollection() users.UserCollection {
	return d.userCollection_
}

func (d databaser) CreateUser(user models.User) (interface{}, error) {
	return d.userCollection().Create(user)
}

func (d databaser) ReadAllUsers() ([]models.User, error) {
	return d.userCollection().ReadAll()
}

func (d databaser) ReadUserById(id primitive.ObjectID) (models.User, error) {
	return d.userCollection().ReadById(id)
}

func (d databaser) UpdateUser(id primitive.ObjectID, user models.User) error {
	return d.userCollection().Update(id, user)
}

func (d databaser) DeleteUserById(id primitive.ObjectID) error {
	return d.userCollection().DeleteById(id)
}

func (d databaser) DeleteAllUsers() error {
	return d.userCollection().DeleteAll()
}
