package assignments

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"verottaa/constants"
	"verottaa/models"
	"verottaa/models/assignments"
	"verottaa/utils"
)

type assignmentController struct {
	collectionName string
	databaseFunc   func(ctx context.Context) *mongo.Database
}

type AssignmentCollection interface {
	models.Destroyable
	getCollection(ctx context.Context) *mongo.Collection
	Creatable
	Readable
	Updatable
	Deletable
}

type Creatable interface {
	Create(assignments.Assignment) (interface{}, error)
	CreateByUserAndPlanIds(primitive.ObjectID, primitive.ObjectID) (interface{}, error)
}

type Readable interface {
	ReadAll() ([]assignments.Assignment, error)
	ReadById(primitive.ObjectID) (assignments.Assignment, error)
}

type Updatable interface {
	Update(primitive.ObjectID, assignments.Assignment) error
}

type Deletable interface {
	DeleteAll() error
	DeleteById(primitive.ObjectID) error
}

var destroyCh = make(chan bool)

var instance *assignmentController
var once sync.Once

func createInstance(databaseFunc func(ctx context.Context) *mongo.Database) *assignmentController {
	inst := new(assignmentController)
	inst.collectionName = constants.AssignmentsCollection
	inst.databaseFunc = databaseFunc
	return inst
}

func GetPlanCollection(databaseFunc func(ctx context.Context) *mongo.Database) AssignmentCollection {
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

func (c assignmentController) Destroy() {
	destroyCh <- true
	close(destroyCh)
	instance = nil
}

func (c assignmentController) getCollection(ctx context.Context) *mongo.Collection {
	collection := c.databaseFunc(ctx).Collection(c.collectionName)
	return collection
}

func (c assignmentController) Create(assignment assignments.Assignment) (interface{}, error) {
	ctx := utils.GetContext()

	insert := assignment.ToBson()

	res, err := c.getCollection(ctx).InsertOne(ctx, insert)
	if err != nil {
		// TODO: логирование
		return nil, err
	}
	return res.InsertedID, nil
}

func (c assignmentController) CreateByUserAndPlanIds(userId primitive.ObjectID, planId primitive.ObjectID) (interface{}, error) {
	assignment := assignments.NewAssignment(userId, planId)
	return c.Create(*assignment)
}

func (c assignmentController) ReadAll() ([]assignments.Assignment, error) {
	ctx := utils.GetContext()
	var allAssignments []assignments.Assignment

	cursor, err := c.getCollection(ctx).Find(ctx, bson.D{})
	if err != nil {
		return allAssignments, err
	}

	for cursor.Next(ctx) {
		var assignment assignments.Assignment
		err := cursor.Decode(&assignment)
		if err != nil {
			return allAssignments, err
		}

		allAssignments = append(allAssignments, assignment)
	}

	if err := cursor.Err(); err != nil {
		return allAssignments, nil
	}

	err = cursor.Close(ctx)

	return allAssignments, err
}

func (c assignmentController) ReadById(id primitive.ObjectID) (assignments.Assignment, error) {
	ctx := utils.GetContext()
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}
	var assignment assignments.Assignment

	err := c.getCollection(ctx).FindOne(ctx, filter).Decode(&assignment)
	if err != nil {
		return assignment, err
	}

	return assignment, err
}

func (c assignmentController) Update(id primitive.ObjectID, assignment assignments.Assignment) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"user_id", assignment.UserId},
			{"plan_d", assignment.PlanId},
			{"curator_id", assignment.CuratorId},
			{"planned_start_date", assignment.PlannedStartDate},
			{"planned_end_date", assignment.PlannedEndDate},
			{"fact_start_date", assignment.FactStartDate},
			{"fact_end_date", assignment.FactEndDate},
			{"current_step_id", assignment.CurrentStepId},
		}},
	}

	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).UpdateOne(ctx, filter, update)
	return err
}

func (c assignmentController) DeleteAll() error {
	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).DeleteMany(ctx, bson.D{})
	return err
}

func (c assignmentController) DeleteById(id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	ctx := utils.GetContext()
	_, err := c.getCollection(ctx).DeleteOne(ctx, filter)
	return err
}
