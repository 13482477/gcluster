package http

import (
	"mcloud/ecs.v2/model"
)

type ListInstancesInput struct {
	AccountId int    `json:"account_id"`
	RegionId  string `json:"region_id"`
	GroupId   int    `json:"group_id"`
}


type EcsModelWithCost struct {
	model.Ecs
	CurrentMonthCost  float64
	PreviousMonthCost float64
}
