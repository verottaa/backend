package assignments

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Assignment struct {
	Id               primitive.ObjectID `json:"id" bson:"_id"`
	UserId           primitive.ObjectID `json:"user_id" bson:"user_id"`
	PlanId           primitive.ObjectID `json:"plan_id" json:"plan_id"`
	CuratorId        primitive.ObjectID `json:"curator_id" bson:"curator_id"`
	PlannedStartDate string             `json:"planned_start_date" bson:"planned_start_date"`
	PlannedEndDate   string             `json:"planned_end_date" bson:"planned_end_date"`
	FactStartDate    string             `json:"fact_start_date" bson:"fact_start_date"`
	FactEndDate      string             `json:"fact_end_date" bson:"fact_end_date"`
	CurrentStepId    primitive.ObjectID `json:"current_step_id" bson:"current_step_id"`
}

func NewAssignment(userId primitive.ObjectID, planId primitive.ObjectID) *Assignment {
	return &Assignment{UserId: userId, PlanId: planId}
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
