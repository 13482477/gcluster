package http

type EcsCreateInput struct {
	AccountId          int    `json:"account_id"`
	AccountName        string `valid:"required" json:"account_name"`
	Region             string `valid:"required" json:"region"`
	ImageId            string `valid:"required" json:"image_id"`
	InstanceName       string `valid:"required" json:"instance_name"`
	VSwitchId          string `valid:"required" json:"vswitch_id"`
	SecurityGroupId    string `valid:"required" json:"security_group_id"`
	IType              string `valid:"required" json:"itype"`
	Uuid               string
	ZoneId             string `valid:"required" json:"zone_id"`
	InstanceChargeType string `valid:"required" json:"instance_charge_type"`
	SystemDiskSize     int    `valid:"required" json:"system_disk_size"`
}
