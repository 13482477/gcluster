package http

type DashboardInput struct {
	AccountId int `json:"account_id"`
	GroupId   int `json:"group_id"`
}

type Dashboard struct {
	CpuMemList       []*EcsCount `json:"cpu_mem_list"`
	InstanceTypeList []*EcsCount `json:"instance_type_list"`
}

type EcsCount struct {
	Item  string `json:"item"`
	Count int    `json:"count"`
}

