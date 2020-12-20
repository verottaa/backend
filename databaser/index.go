package databaser

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"verottaa/config"
	"verottaa/constants"
	"verottaa/databaser/assignments"
	"verottaa/databaser/plans"
	"verottaa/models"
	assignmentModel "verottaa/models/assignments"
	plansModel "verottaa/models/plans"
	"verottaa/utils"
)

var configuration = config.GetConfiguration()

type databaser struct {
	client                *mongo.Client
	planCollection_       plans.PlanCollection
	assignmentCollection_ assignments.AssignmentCollection
}

type DB interface {
	models.Destroyable
	PlansCollection
	AssignmentsCollection
}

type PlansCollection interface {
	CreatePlan(plansModel.Plan) (interface{}, error)
	ReadAllPlans() ([]plansModel.Plan, error)
	ReadPlanById(primitive.ObjectID) (plansModel.Plan, error)
	UpdatePlan(primitive.ObjectID, plansModel.Plan) error
	DeletePlanById(primitive.ObjectID) error
	DeleteAllPlans() error

	CreateStepInPlan(id primitive.ObjectID, step plansModel.Step) (interface{}, error)
	ReadAllStepsInPlan(id primitive.ObjectID) ([]plansModel.Step, error)
	ReadStepByIdInPlan(planId primitive.ObjectID, stepId primitive.ObjectID) (plansModel.Step, error)
	UpdateStepInPlan(id primitive.ObjectID, stepId primitive.ObjectID, updateStep plansModel.Step) error
	DeleteStepInPlan(planId primitive.ObjectID, stepId primitive.ObjectID) error
	DeleteAllStepsInPlan(id primitive.ObjectID) error
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
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "initDatabaser",
			"error":    err,
			"cause":    "initialisation new mongo client",
		}).Error("Unexpected error")
	}

	db.planCollection_ = plans.GetPlanCollection(database)
	db.assignmentCollection_ = assignments.GetPlanCollection(database)

	return db
}

func database(ctx context.Context) *mongo.Database {
	err := instance.client.Connect(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "database",
			"error":    err,
			"cause":    "Trying to connect with client to mongo",
		}).Error("Unexpected error")
	}

	err = instance.client.Ping(ctx, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "database",
			"error":    err,
			"cause":    "Trying to ping mongo client",
		}).Error("Unexpected error")
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
	instance.planCollection().Destroy()
	instance = nil
}

func (d databaser) planCollection() plans.PlanCollection {
	return d.planCollection_
}

func (d databaser) assignmentCollection() assignments.AssignmentCollection {
	return d.assignmentCollection_
}

//
//	PLANS:
//

func (d databaser) CreatePlan(plan plansModel.Plan) (interface{}, error) {
	return d.planCollection().Create(plan)
}

func (d databaser) ReadAllPlans() ([]plansModel.Plan, error) {
	return d.planCollection().ReadAll()
}

func (d databaser) ReadPlanById(id primitive.ObjectID) (plansModel.Plan, error) {
	return d.planCollection().ReadById(id)
}

func (d databaser) UpdatePlan(id primitive.ObjectID, plan plansModel.Plan) error {
	return d.planCollection().Update(id, plan)
}

func (d databaser) DeletePlanById(id primitive.ObjectID) error {
	return d.planCollection().DeleteById(id)
}

func (d databaser) DeleteAllPlans() error {
	return d.planCollection().DeleteAll()
}

//
//	plans/STEPS
//

func (d databaser) CreateStepInPlan(planId primitive.ObjectID, step plansModel.Step) (interface{}, error) {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "CreateStepInPlan",
			"error":    err,
			"cause":    "read plan by id",
		}).Error("Unexpected error")
		return nil, err
	}
	var stepId = utils.NewObjectId()
	step.Id = stepId
	plan.AddStep(step)
	err = d.UpdatePlan(planId, plan)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "CreateStepInPlan",
			"error":    err,
			"cause":    "updating plan",
		}).Error("Unexpected error")
		return nil, err
	}
	return stepId, err
}

func (d databaser) ReadAllStepsInPlan(planId primitive.ObjectID) ([]plansModel.Step, error) {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "ReadAllStepsInPlan",
			"error":    err,
			"cause":    "read plan by id",
		}).Error("Unexpected error")
		return nil, err
	}
	return plan.Steps, nil
}

func (d databaser) ReadStepByIdInPlan(planId primitive.ObjectID, stepId primitive.ObjectID) (plansModel.Step, error) {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "ReadStepByIdInPlan",
			"error":    err,
			"cause":    "read plan by id",
		}).Error("Unexpected error")
		return plansModel.Step{}, err
	}
	step, err := plan.GetStepById(stepId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "ReadStepByIdInPlan",
			"error":    err,
			"cause":    "getting step from plan by id",
		}).Error("Unexpected error")
		return plansModel.Step{}, err
	}
	return step, nil
}

func (d databaser) UpdateStepInPlan(planId primitive.ObjectID, stepId primitive.ObjectID, updatedStep plansModel.Step) error {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "UpdateStepInPlan",
			"error":    err,
			"cause":    "read plan by id",
		}).Error("Unexpected error")
		return err
	}
	err = plan.UpdateStep(stepId, updatedStep)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "UpdateStepInPlan",
			"error":    err,
			"cause":    "update step",
		}).Error("Unexpected error")
		return err
	}
	err = d.UpdatePlan(planId, plan)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "UpdateStepInPlan",
			"error":    err,
			"cause":    "updating plan",
		}).Error("Unexpected error")
		return err
	}
	return nil
}

func (d databaser) DeleteStepInPlan(planId primitive.ObjectID, stepId primitive.ObjectID) error {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "DeleteStepInPlan",
			"error":    err,
			"cause":    "reading plan by id",
		}).Error("Unexpected error")
		return err
	}
	err = plan.RemoveStep(stepId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "DeleteStepInPlan",
			"error":    err,
			"cause":    "removing step from plan",
		}).Error("Unexpected error")
		return err
	}
	err = d.UpdatePlan(planId, plan)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "DeleteStepInPlan",
			"error":    err,
			"cause":    "updating plan",
		}).Error("Unexpected error")
		return err
	}
	return nil
}

func (d databaser) DeleteAllStepsInPlan(planId primitive.ObjectID) error {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "DeleteAllStepsInPlan",
			"error":    err,
			"cause":    "read plan by id",
		}).Error("Unexpected error")
		return err
	}
	plan.RemoveAllSteps()
	err = d.UpdatePlan(planId, plan)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "databaser",
			"function": "DeleteAllStepsInPlan",
			"error":    err,
			"cause":    "updating plan",
		}).Error("Unexpected error")
		return err
	}
	return nil
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
