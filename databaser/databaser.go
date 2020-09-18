package databaser

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"verottaa/models"
	"verottaa/variables"
)

type UserCollection struct {
	collection *mongo.Collection
}

var User = new(UserCollection)

func (c *UserCollection) Init() {
	c.collection = variables.Client.Database(variables.DATABASE_NAME).Collection(variables.USERS_COLLECTION)
}

func (c *UserCollection) Create(user models.User) error {
	_, err := c.collection.InsertOne(variables.Ctx, user)
	return err
}

func (c *UserCollection) Read() ([]models.User, error) {
	var users []models.User

	cursor, err := c.collection.Find(variables.Ctx, bson.D{})
	if err != nil {
		return users, err
	}

	for cursor.Next(variables.Ctx) {
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

	err = cursor.Close(variables.Ctx)

	return users, err
}

func (c *UserCollection) ReadById(id primitive.ObjectID) (models.User, error) {
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}
	var user models.User

	err := c.collection.FindOne(variables.Ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, err
}

func (c *UserCollection) Update(id primitive.ObjectID, user models.User) error {
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

	_, err := c.collection.UpdateOne(variables.Ctx, filter, update)
	return err
}

func (c *UserCollection) DeleteById(id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := c.collection.DeleteOne(variables.Ctx, filter)
	return err
}

func (c *UserCollection) Delete() error {
	_, err := c.collection.DeleteMany(variables.Ctx, bson.D{})
	return err
}
