package plans

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"verottaa/constants"
	"verottaa/models"
	"verottaa/models/plans"
	"verottaa/utils"
)

type planController struct {
	collectionName string
	databaseFunc   func(ctx context.Context) *mongo.Database
}

type PlanCollection interface {
	models.Destroyable
	getCollection(ctx context.Context) *mongo.Collection
	Creatable
	Readable
	Updatable
	Deletable
}

type Creatable interface {
	Create(plans.Plan) (interface{}, error)
}

type Readable interface {
	ReadAll() ([]plans.Plan, error)
	ReadById(primitive.ObjectID) (plans.Plan, error)
}

type Updatable interface {
	Update(primitive.ObjectID, plans.Plan) error
}

type Deletable interface {
	DeleteAll() error
	DeleteById(primitive.ObjectID) error
}

var destroyCh = make(chan bool)

var instance *planController
var once sync.Once

func createInstance(databaseFunc func(ctx context.Context) *mongo.Database) *planController {
	inst := new(planController)
	inst.collectionName = constants.PlansCollection
	inst.databaseFunc = databaseFunc
	return inst
}

func GetPlanCollection(databaseFunc func(ctx context.Context) *mongo.Database) PlanCollection {
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

func (c planController) Destroy() {
	destroyCh <- true
	close(destroyCh)
	instance = nil
}

func (c planController) getCollection(ctx context.Context) *mongo.Collection {
	collection := c.databaseFunc(ctx).Collection(c.collectionName)
	return collection
}

func (c planController) Create(plan plans.Plan) (interface{}, error) {
	ctx := utils.GetContext()

	plan.RecalculatePeriod()
	insert := plan.ToBson()

	res, err := c.getCollection(ctx).InsertOne(ctx, insert)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "Create",
			"error":    err,
			"cause":    "inserting plan to database",
		}).Error("Unexpected error")
		return nil, err
	}
	return res.InsertedID, nil
}

func (c planController) ReadAll() ([]plans.Plan, error) {
	ctx := utils.GetContext()
	var allPlans []plans.Plan

	cursor, err := c.getCollection(ctx).Find(ctx, bson.D{})
	if err != nil {
		return allPlans, err
	}

	for cursor.Next(ctx) {
		var plan plans.Plan
		err := cursor.Decode(&plan)
		if err != nil {
			return allPlans, err
		}

		allPlans = append(allPlans, plan)
	}

	if err := cursor.Err(); err != nil {
		return allPlans, nil
	}

	err = cursor.Close(ctx)

	return allPlans, err
}

func (c planController) ReadById(id primitive.ObjectID) (plans.Plan, error) {
	ctx := utils.GetContext()
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}
	var plan plans.Plan

	err := c.getCollection(ctx).FindOne(ctx, filter).Decode(&plan)
	if err != nil {
		return plan, err
	}

	return plan, err
}

func (c planController) Update(id primitive.ObjectID, plan plans.Plan) error {
	// Пересчитываем период для того,
	// чтобы это было корректное число,
	// потому что на фронте это выполняться не будет
	plan.RecalculatePeriod()
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"steps", plan.Steps},
			{"period", plan.Period},
		}},
	}

	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).UpdateOne(ctx, filter, update)
	return err
}

func (c planController) DeleteAll() error {
	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).DeleteMany(ctx, bson.D{})
	return err
}

func (c planController) DeleteById(id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).DeleteOne(ctx, filter)
	return err
}
