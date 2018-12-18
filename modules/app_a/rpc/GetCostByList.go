package rpc

type CostList struct {
	InstanceId string  `json:"instance_id"`
	Cost       float64 `json:"cost"`
}

type GetCostByListRequest struct {
	InstanceList []string `json:"instance_list"`
	BillingCycle string   `json:"billing_cycle"`
}

type GetCostByListResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    []CostList `json:"data"`
}
