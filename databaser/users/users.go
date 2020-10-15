package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"verottaa/constants"
	"verottaa/models"
	"verottaa/models/users"
	"verottaa/utils"
)

type userController struct {
	collectionName string
	databaseFunc   func(ctx context.Context) *mongo.Database
}

type UserCollection interface {
	models.Destroyable
	getCollection(ctx context.Context) *mongo.Collection
	Creatable
	Readable
	Updatable
	Deletable
}

type Creatable interface {
	Create(user users.User) (interface{}, error)
}

type Readable interface {
	ReadAll() ([]users.User, error)
	ReadById(id primitive.ObjectID) (users.User, error)
}

type Updatable interface {
	Update(id primitive.ObjectID, user users.User) error
}

type Deletable interface {
	DeleteAll() error
	DeleteById(id primitive.ObjectID) error
}

var destroyCh = make(chan bool)

var instance *userController
var once sync.Once

func createInstance(databaseFunc func(ctx context.Context) *mongo.Database) *userController {
	inst := new(userController)
	inst.collectionName = constants.UsersCollection
	inst.databaseFunc = databaseFunc
	return inst
}

func GetUserCollection(databaseFunc func(ctx context.Context) *mongo.Database) UserCollection {
	once.Do(func() {
		instance = createInstance(databaseFunc)
		go func() {
			for {
				select {
				case <-destroyCh:
					return
				}
			}
		}()
	})

	return instance
}

func (c userController) Destroy() {
	destroyCh <- true
	close(destroyCh)
	instance = nil
}

func (c userController) getCollection(ctx context.Context) *mongo.Collection {
	collection := c.databaseFunc(ctx).Collection(c.collectionName)
	return collection
}

func (c userController) Create(user users.User) (interface{}, error) {
	ctx := utils.GetContext()

	insert := user.ToBson()

	res, err := c.getCollection(ctx).InsertOne(ctx, insert)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, err
}

func (c userController) ReadAll() ([]users.User, error) {
	ctx := utils.GetContext()
	var allUsers []users.User

	cursor, err := c.getCollection(ctx).Find(ctx, bson.D{})
	if err != nil {
		return allUsers, err
	}

	for cursor.Next(ctx) {
		var user users.User
		err := cursor.Decode(&user)
		if err != nil {
			return allUsers, err
		}

		allUsers = append(allUsers, user)
	}

	if err := cursor.Err(); err != nil {
		return allUsers, nil
	}

	err = cursor.Close(ctx)

	return allUsers, err
}

func (c userController) ReadById(id primitive.ObjectID) (users.User, error) {
	ctx := utils.GetContext()
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}
	var user users.User

	err := c.getCollection(ctx).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, err
}

func (c userController) Update(id primitive.ObjectID, user users.User) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"firstName", user.FirstName},
			{"secondName", user.SecondName},
			{"patronymic", user.Patronymic},
			{"type", user.Type},
			{"branch", user.Branch},
			{"department", user.Department},
		}},
	}

	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).UpdateOne(ctx, filter, update)
	return err
}

func (c userController) DeleteAll() error {
	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).DeleteMany(ctx, bson.D{})
	return err
}

func (c userController) DeleteById(id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).DeleteOne(ctx, filter)
	return err
}
