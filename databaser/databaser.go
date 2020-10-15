package databaser

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"verottaa/config"
	"verottaa/constants"
	"verottaa/databaser/plans"
	"verottaa/databaser/users"
	"verottaa/models"
	"verottaa/utils"
	logpack "verottaa/utils/logger"
)

var configuration = config.GetConfiguration()

type databaser struct {
	client          *mongo.Client
	userCollection_ users.UserCollection
	planCollection_ plans.PlanCollection
}

type DB interface {
	models.Destroyable
	UserCollection
	PlansCollection
}

type UserCollection interface {
	CreateUser(models.User) (interface{}, error)
	ReadAllUsers() ([]models.User, error)
	ReadUserById(primitive.ObjectID) (models.User, error)
	UpdateUser(primitive.ObjectID, models.User) error
	DeleteUserById(primitive.ObjectID) error
	DeleteAllUsers() error
}

type PlansCollection interface {
	CreatePlan(models.Plan) (interface{}, error)
	ReadAllPlans() ([]models.Plan, error)
	ReadPlanById(primitive.ObjectID) (models.Plan, error)
	UpdatePlan(primitive.ObjectID, models.Plan) error
	DeletePlanById(primitive.ObjectID) error
	DeleteAllPlans() error

	CreateStepInPlan(id primitive.ObjectID, step models.Step) (interface{}, error)
	ReadAllStepsInPlan(id primitive.ObjectID) ([]models.Step, error)
	ReadStepByIdInPlan(planId primitive.ObjectID, stepId primitive.ObjectID) (models.Step, error)
	UpdateStepInPlan(id primitive.ObjectID, stepId primitive.ObjectID, updateStep models.Step) error
	DeleteStepInPlan(planId primitive.ObjectID, stepId primitive.ObjectID) error
	DeleteAllStepsInPlan(id primitive.ObjectID) error
}

var destroyCh = make(chan bool)

var instance *databaser
var once sync.Once

var logTag = "DATABASER"
var logger *logpack.Logger

func init() {
	logger = logpack.CreateLogger(logTag)
}

func initDatabaser() *databaser {
	db := new(databaser)

	var err error
	db.client, err = mongo.NewClient(options.Client().ApplyURI(configuration.GetDatabaseHost()))
	if err != nil {
		logger.Error(err)
	}

	db.userCollection_ = users.GetUserCollection(database)
	db.planCollection_ = plans.GetPlanCollection(database)

	return db
}

func database(ctx context.Context) *mongo.Database {
	err := instance.client.Connect(ctx)
	if err != nil {
		logger.Error(err)
	}

	err = instance.client.Ping(ctx, nil)
	if err != nil {
		logger.Error(err)
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
	instance.planCollection().Destroy()
	instance = nil
}

func (d databaser) userCollection() users.UserCollection {
	return d.userCollection_
}

func (d databaser) planCollection() plans.PlanCollection {
	return d.planCollection_
}

//
//	USERS
//

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

//
//	PLANS:
//

func (d databaser) CreatePlan(plan models.Plan) (interface{}, error) {
	return d.planCollection().Create(plan)
}

func (d databaser) ReadAllPlans() ([]models.Plan, error) {
	return d.planCollection().ReadAll()
}

func (d databaser) ReadPlanById(id primitive.ObjectID) (models.Plan, error) {
	return d.planCollection().ReadById(id)
}

func (d databaser) UpdatePlan(id primitive.ObjectID, plan models.Plan) error {
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

func (d databaser) CreateStepInPlan(planId primitive.ObjectID, step models.Step) (interface{}, error) {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		// TODO: logger
		return nil, err
	}
	var stepId = utils.NewObjectId()
	step.Id = stepId
	plan.AddStep(step)
	err = d.UpdatePlan(planId, plan)
	if err != nil {
		//TODO: logger
		return nil, err
	}
	return stepId, err
}

func (d databaser) ReadAllStepsInPlan(planId primitive.ObjectID) ([]models.Step, error) {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		// TODO: logger
		return nil, err
	}
	return plan.Steps, nil
}

func (d databaser) ReadStepByIdInPlan(planId primitive.ObjectID, stepId primitive.ObjectID) (models.Step, error) {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		// TODO: logger
		return models.Step{}, err
	}
	step, err := plan.GetStepById(stepId)
	if err != nil {
		// TODO: logger
		return models.Step{}, err
	}
	return step, nil
}

func (d databaser) UpdateStepInPlan(planId primitive.ObjectID, stepId primitive.ObjectID, updatedStep models.Step) error {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		// TODO: logger
		return err
	}
	err = plan.UpdateStep(stepId, updatedStep)
	if err != nil {
		// TODO: logger
		return err
	}
	err = d.UpdatePlan(planId, plan)
	if err != nil {
		//TODO: logger
		return err
	}
	return nil
}

func (d databaser) DeleteStepInPlan(planId primitive.ObjectID, stepId primitive.ObjectID) error {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		// TODO: logger
		return err
	}
	err = plan.RemoveStep(stepId)
	if err != nil {
		// TODO: logger
		return err
	}
	err = d.UpdatePlan(planId, plan)
	if err != nil {
		//TODO: logger
		return err
	}
	return nil
}

func (d databaser) DeleteAllStepsInPlan(planId primitive.ObjectID) error {
	plan, err := d.ReadPlanById(planId)
	if err != nil {
		// TODO: logger
		return err
	}
	plan.RemoveAllSteps()
	err = d.UpdatePlan(planId, plan)
	if err != nil {
		//TODO: logger
		return err
	}
	return nil
}
