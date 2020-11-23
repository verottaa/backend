package assignments

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"verottaa/controllers/controller_helpers"
	"verottaa/databaser"
	"verottaa/models/dto"
	"verottaa/models/plans"
	"verottaa/utils"
)

var db = databaser.GetDatabaser()
var variableReader = controller_helpers.GetVariableReader()

type Assignable interface {
	RecalculateEndDates()
	RecalculatePlannedEndDate(plans.Plan)
	RecalculateFactEndDate(plans.Plan)
	SetCurrentStepWithIndex(int)
}

type Assignment struct {
	Id               primitive.ObjectID `json:"id" bson:"_id"`
	UserId           primitive.ObjectID `json:"user_id" bson:"user_id"`
	PlanId           primitive.ObjectID `json:"plan_id" json:"plan_id"`
	CuratorId        primitive.ObjectID `json:"curator_id" bson:"curator_id"`
	PlannedStartDate time.Time          `json:"planned_start_date" bson:"planned_start_date"`
	PlannedEndDate   time.Time          `json:"planned_end_date" bson:"planned_end_date"`
	FactStartDate    time.Time          `json:"fact_start_date" bson:"fact_start_date"`
	FactEndDate      time.Time          `json:"fact_end_date" bson:"fact_end_date"`
	CurrentStepId    primitive.ObjectID `json:"current_step_id" bson:"current_step_id"`
}

func (a *Assignment) RecalculateEndDates() {

	plan, err := db.ReadPlanById(a.PlanId)
	if err != nil {

	}

	a.RecalculatePlannedEndDate(plan)
	a.RecalculateFactEndDate(plan)
}

func (a *Assignment) RecalculatePlannedEndDate(plan plans.Plan) {
	a.PlannedEndDate = a.PlannedStartDate.Add(time.Hour * 24 * time.Duration(plan.Period))
}

func (a *Assignment) RecalculateFactEndDate(plan plans.Plan) {
	a.FactEndDate = a.FactStartDate.Add(time.Hour * 24 * time.Duration(plan.Period))
}

func (a *Assignment) SetCurrentStepWithIndex(index int) {

	plan, err := db.ReadPlanById(a.PlanId)
	if err != nil {
		a.CurrentStepId = primitive.ObjectID{}
	}

	if len(plan.Steps) != 0 {
		id := plan.Steps[index].Id
		a.CurrentStepId = id
	}
}

func NewAssignment(assignDto dto.AssignCreateDto) Assignment {
	// TODO: сделать назначение куратора
	assignment := Assignment{
		UserId:           variableReader.GetObjectIDFromString(assignDto.UserId),
		PlanId:           variableReader.GetObjectIDFromString(assignDto.PlanId),
		CuratorId:        variableReader.GetObjectIDFromString(assignDto.CuratorId),
		PlannedStartDate: utils.ParseTime(assignDto.StartDate),
		PlannedEndDate:   nil,
		FactStartDate:    utils.ParseTime(assignDto.StartDate),
		FactEndDate:      nil,
		CurrentStepId:    primitive.ObjectID{},
	}

	assignment.RecalculateEndDates()
	assignment.SetCurrentStepWithIndex(0)

	return assignment
}

func (a *Assignment) ToBson() bson.M {
	return bson.M{
		"user_id":            a.UserId,
		"plan_id":            a.PlanId,
		"curator_id":         a.CuratorId,
		"planned_start_date": a.PlannedStartDate,
		"planned_end_date":   a.PlannedEndDate,
		"fact_start_date":    a.FactStartDate,
		"fact_end_date":      a.FactEndDate,
		"current_step_id":    a.CurrentStepId,
	}
}
