package databaser

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"verottaa/config"
	"verottaa/constants"
	"verottaa/databaser/assignments"
	"verottaa/models"
	assignmentModel "verottaa/models/assignments"
)

var configuration = config.GetConfiguration()

type databaser struct {
	client                *mongo.Client
	assignmentCollection_ assignments.AssignmentCollection
}

type DB interface {
	models.Destroyable
	AssignmentsCollection
}

type AssignmentsCollection interface {
	CreateAssignment(assignmentModel.Assignment) (interface{}, error)
	ReadAllAssignments() ([]assignmentModel.Assignment, error)
	ReadAssignmentById(primitive.ObjectID) (assignmentModel.Assignment, error)
	UpdateAssignment(primitive.ObjectID, assignmentModel.Assignment) error
	DeleteAssignmentById(primitive.ObjectID) error
	DeleteAllAssignments() error
}

var destroyCh = make(chan bool)

var instance *databaser
var once sync.Once

func initDatabaser() *databaser {
	db := new(databaser)

	var err error
	db.client, err = mongo.NewClient(options.Client().ApplyURI(configuration.GetDatabaseHost()))
	if err != nil {

	}

	db.assignmentCollection_ = assignments.GetPlanCollection(database)

	return db
}

func database(ctx context.Context) *mongo.Database {
	err := instance.client.Connect(ctx)
	if err != nil {

	}

	err = instance.client.Ping(ctx, nil)
	if err != nil {

	}
	return instance.client.Database(constants.DatabaseName)
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
	instance = nil
}

func (d databaser) assignmentCollection() assignments.AssignmentCollection {
	return d.assignmentCollection_
}

//
//	ASSIGNMENTS
//

func (d databaser) CreateAssignment(assignment assignmentModel.Assignment) (interface{}, error) {
	return d.assignmentCollection().Create(assignment)
}

func (d databaser) ReadAllAssignments() ([]assignmentModel.Assignment, error) {
	return d.assignmentCollection().ReadAll()
}

func (d databaser) ReadAssignmentById(id primitive.ObjectID) (assignmentModel.Assignment, error) {
	return d.assignmentCollection().ReadById(id)
}

func (d databaser) UpdateAssignment(id primitive.ObjectID, assignment assignmentModel.Assignment) error {
	return d.assignmentCollection().Update(id, assignment)
}

func (d databaser) DeleteAssignmentById(id primitive.ObjectID) error {
	return d.assignmentCollection().DeleteById(id)
}

func (d databaser) DeleteAllAssignments() error {
	return d.assignmentCollection().DeleteAll()
}
