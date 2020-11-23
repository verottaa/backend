package dto

type ObjectCreatedDto struct {
	Id string `json:"id"`
}

type AssignCreateDto struct {
	PlanId    string `json:"planId"`
	UserId    string `json:"userId"`
	StartDate string `json:"startDate"`
	CuratorId string `json:"curatorId"`
}
