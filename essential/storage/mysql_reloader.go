package storage

import (
	"time"

	"github.com/go-xorm/xorm"
)

type MysqlReloader struct {
	DB       *xorm.Engine
	Ticker   *time.Ticker
	Reloader *DataReloader
}

type DataReloader interface {
	Query(reloader *MysqlReloader) error
}

func NewMysqlReloader(db *xorm.Engine, reloader DataReloader, intervalMSec int) (*MysqlReloader, error) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(intervalMSec))

	result := &MysqlReloader{
		DB:     db,
		Ticker: ticker,
	}

	reloader.Query(result)

	go func() {
		for range ticker.C {
			reloader.Query(result)
		}
	}()

	return result, nil
}
