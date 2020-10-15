package models

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	SecondName string             `json:"secondName" bson:"secondName"`
	Patronymic string             `json:"patronymic" bson:"patronymic"`
	Type       string             `json:"type" bson:"type"`
	Branch     string             `json:"branch" bson:"branch"`
	Department string             `json:"department" bson:"department"`
}

type Step struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Title         string             `json:"title" bson:"title"`
	EducationForm string             `json:"education_form" bson:"education_form"`
	Period        int                `json:"period" bson:"period"`
	Materials     string             `json:"materials" bson:"materials"`
}

type Plan struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Steps  []Step             `json:"steps" bson:"steps"`
	Period int                `json:"period" bson:"period"`
}

func (p *Plan) AddStep(step Step) {
	p.Steps = append(p.Steps, step)
	p.RecalculatePeriod()
}

func (p *Plan) RemoveStep(id primitive.ObjectID) error {
	var index, _, err = p.findStepIndexById(id)
	if err != nil {
		// TODO: логирование
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
		// TODO: логирование
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
		// TODO: логирование
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
