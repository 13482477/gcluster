package http

type EcsDescribeInstanceAttrInput struct {
	Platform   int    `json:"platform"`
	AccountId  int    `json:"account_id"`
	Region     string `json:"region"`
	InstanceId string `json:"instance_id"`
}

