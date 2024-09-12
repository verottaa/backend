package entity

import (
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Plan struct {
	mgm.DefaultModel `bson:",inline"`
	Steps  []Step `json:"steps" bson:"steps"`
	Period int    `json:"period" bson:"period"`
}

type PlanFilters struct {
	//Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Steps  []Step `json:"steps,omitempty" bson:"steps,omitempty"`
	Period int    `json:"period,omitempty" bson:"period,omitempty"`
}

type Step struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Title         string             `json:"title" bson:"title"`
	EducationForm string             `json:"education_form" bson:"education_form"`
	Period        int                `json:"period" bson:"period"`
	Materials     string             `json:"materials" bson:"materials"`
}

type StepFilters struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title         string             `json:"title,omitempty" bson:"title,omitempty"`
	EducationForm string             `json:"education_form,omitempty" bson:"education_form,omitempty"`
	Period        int                `json:"period,omitempty" bson:"period,omitempty"`
	Materials     string             `json:"materials,omitempty" bson:"materials,omitempty"`
}

// TODO: вынести методы в сервис

func (p *Plan) AddStep(step Step) {
	p.Steps = append(p.Steps, step)
	p.RecalculatePeriod()
}

func (p *Plan) RemoveStep(id primitive.ObjectID) error {
	var index, _, err = p.findStepIndexById(id)
	if err != nil {
		return err
	}
	p.Steps = append(p.Steps[:index], p.Steps[index+1:]...)
	p.RecalculatePeriod()
	return nil
}

func (p *Plan) RemoveAllSteps() {
	p.Steps = make([]Step, 0)
	p.RecalculatePeriod()
}

func (p *Plan) UpdateStep(id primitive.ObjectID, step Step) error {
	var index, _, err = p.findStepIndexById(id)
	if err != nil {
		return err
	}
	step.Id = id
	p.Steps[index] = step
	p.RecalculatePeriod()
	return nil
}

func (p *Plan) GetSteps() []Step {
	return p.Steps
}

func (p *Plan) GetStepById(id primitive.ObjectID) (Step, error) {
	var _, step, err = p.findStepIndexById(id)
	if err != nil {
		return step, err
	}
	return step, nil
}

func (p *Plan) findStepIndexById(id primitive.ObjectID) (int, Step, error) {
	for i, step := range p.Steps {
		if step.Id == id {
			return i, step, nil
		}
	}
	return 0, Step{}, errors.New("can't find index of element")
}

func (p *Plan) RecalculatePeriod() {
	var period = 0
	for _, step := range p.Steps {
		period = period + step.Period
	}
	p.Period = period
}
