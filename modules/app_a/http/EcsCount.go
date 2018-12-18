package http

type EcsCountRequest struct {
	AccountId int    `json:"account_id"`
	RegionId  string `json:"region_id"`
	GroupId   int    `json:"group_id"`
}

type Count struct {
	Count int `json:"count"`
}

