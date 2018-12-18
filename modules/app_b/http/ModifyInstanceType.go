package http

type ModifyInstanceTypeInput struct {
	AccountId    int    `valid:"required" json:"account_id"`
	Region       string `valid:"required" json:"region"`
	InstanceId   string `valid:"required" json:"instance_id"`
	Uuid         string
	InstanceType string `valid:"required" json:"instance_type"`
}
