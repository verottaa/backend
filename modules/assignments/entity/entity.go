package entity

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Assignment struct {
	mgm.DefaultModel `bson:",inline"`
	UserId           primitive.ObjectID `json:"user_id" bson:"user_id"`
	PlanId           primitive.ObjectID `json:"plan_id" json:"plan_id"`
	CuratorId        primitive.ObjectID `json:"curator_id" bson:"curator_id"`
	PlannedStartDate time.Time          `json:"planned_start_date" bson:"planned_start_date"`
	PlannedEndDate   time.Time          `json:"planned_end_date" bson:"planned_end_date"`
	FactStartDate    time.Time          `json:"fact_start_date" bson:"fact_start_date"`
	FactEndDate      time.Time          `json:"fact_end_date" bson:"fact_end_date"`
	CurrentStepId    primitive.ObjectID `json:"current_step_id" bson:"current_step_id"`
}

type AssignCreateDto struct {
	PlanId    string `json:"planId"`
	UserId    string `json:"userId"`
	StartDate string `json:"startDate"`
	CuratorId string `json:"curatorId"`
}
