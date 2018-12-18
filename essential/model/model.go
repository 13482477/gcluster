package model

import (
	"time"
)

type Platform int

const (
	_                       Platform = iota
	PlatformAliCloud
	PlatformTencentCloud
	PlatformPhysicalMachine
)

type BaseModel struct {
	ID          int64     `gorm:"primary_key" json:"id"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
	Ext         string    `json:"ext"`
}
