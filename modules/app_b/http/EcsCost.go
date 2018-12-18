package http

type EcsCostRequest struct {
	AccountId int    `json:"account_id"`
	RegionId  string `json:"region_id"`
	GroupId   int    `json:"group_id"`
}

type Cost struct {
	Cost float64 `json:"cost"`
}

