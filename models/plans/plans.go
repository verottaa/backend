package plans

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Step struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Title         string             `json:"title" bson:"title"`
	EducationForm string             `json:"education_form" bson:"education_form"`
	Period        int                `json:"period" bson:"period"`
	Materials     string             `json:"materials" bson:"materials"`
}

func (s *Step) ToBson() bson.M {
	return bson.M{
		"title":          s.Title,
		"education_form": s.EducationForm,
		"period":         s.Period,
		"materials":      s.Materials,
	}
}

type Plan struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Steps  []Step             `json:"steps" bson:"steps"`
	Period int                `json:"period" bson:"period"`
}

func (p *Plan) ToBson() bson.M {
	return bson.M{
		"steps":  p.Steps,
		"period": p.Period,
	}
}

func (p *Plan) AddStep(step Step) {
	p.Steps = append(p.Steps, step)
	p.RecalculatePeriod()
}

func (p *Plan) RemoveStep(id primitive.ObjectID) error {
	var index, _, err = p.findStepIndexById(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "RemoveStep",
			"error":    err,
			"cause":    "Removing step",
		}).Error("Unexpected error")
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
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "UpdateStep",
			"error":    err,
			"cause":    "finding step index",
		}).Error("Unexpected error")
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
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "GetStepById",
			"error":    err,
			"cause":    "finding step by id",
		}).Error("Unexpected error")
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
