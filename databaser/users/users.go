package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"verottaa/constants"
	"verottaa/models"
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
	Create(user models.User) (interface{}, error)
}

type Readable interface {
	ReadAll() ([]models.User, error)
	ReadById(id primitive.ObjectID) (models.User, error)
}

type Updatable interface {
	Update(id primitive.ObjectID, user models.User) error
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
	inst.collectionName = constants.USERS_COLLECTION
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

func (c userController) Create(user models.User) (interface{}, error) {
	ctx := utils.GetContext()

	insert := bson.M{
		"firstName":  user.FirstName,
		"secondName": user.SecondName,
		"patronymic": user.Patronymic,
		"type":       user.Type,
		"branch":     user.Branch,
		"department": user.Department,
	}

	res, err := c.getCollection(ctx).InsertOne(ctx, insert)
	return res.InsertedID, err
}

func (c userController) ReadAll() ([]models.User, error) {
	ctx := utils.GetContext()
	var users []models.User

	cursor, err := c.getCollection(ctx).Find(ctx, bson.D{})
	if err != nil {
		return users, err
	}

	for cursor.Next(ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return users, nil
	}

	err = cursor.Close(ctx)

	return users, err
}

func (c userController) ReadById(id primitive.ObjectID) (models.User, error) {
	ctx := utils.GetContext()
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}
	var user models.User

	err := c.getCollection(ctx).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, err
}

func (c userController) Update(id primitive.ObjectID, user models.User) error {
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
