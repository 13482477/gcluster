package model

import (
	"mcloud/public.v2/model"
)

type Ecs struct {
	model.BaseModel
	Platform     model.Platform `json:"platform" gorm:"unique_index:ecs_index"`
	InstanceId   string         `json:"instance_id" gorm:"unique_index:ecs_index"`
	InstanceName string         `json:"instance_name"`
	PublicIp     string         `json:"public_ip"`  //公网ip
	PrivateIp    string         `json:"private_ip"` //内网ip
	OsType       string         `json:"os_type"`
	Cpu          int            `json:"cpu"`
	Mem          int            `json:"mem"`
	DiskSize     string         `json:"disk_size"`
	InstanceType string         `json:"instance_type"`
	ZoneId       string         `json:"zone_id"`
	RegionId     string         `json:"region_id"`
	GrafanaUrl   string         `json:"grafana_url"`  // grafana监控url
	AccountId    int64          `json:"account_id"`
	AccountName  string         `json:"account_name"`
	Status       EcsStatus      `json:"status" sql:"DEFAULT:0"` // 0 运行中 1 创建中 2 启动中 3 初始化中 -1 停止
	GroupId      int64          `json:"group_id"`
}

type EcsStatus int

const (
	EcsStatusNon          EcsStatus = iota
	EcsStatusCreated
	EcsStatusCreateFailed
	EcsStatusStarting
	EcsStatusInitializing
	EcsStatusRunning
	EcsStatusStopping
	EcsStatusStopped
	EcsStatusReleased
)

type EcsAgent struct {
	model.BaseModel
	EcsId        int64  `json:"ecs_id"`
	PrivateIp    string `json:"private_ip"`
	CheckUrl     string `json:"check_url"`
	CheckUrlSign string `json:"check_url_sign"`
	AgentStatus  int    `json:"agent_status"`
}
