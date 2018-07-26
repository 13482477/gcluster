package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type MysqlStorage struct {
	db *xorm.Engine
}

type MysqlSession struct {
	Session *xorm.Session
}

type MysqlOption struct {
	ConnectionString string
	MaxOpenConn      int
	MaxIdleConn      int
	UseCache         bool
	MaxCacheCount    int
	Expire           int
}

func CreateMysqlStorage(option *MysqlOption) (*MysqlStorage, error) {
	db, err := xorm.NewEngine("mysql", option.ConnectionString)
	if err != nil {
		return nil, err
	}

	if 0 != option.MaxOpenConn {
		db.SetMaxOpenConns(option.MaxOpenConn)
	}

	if 0 != option.MaxIdleConn {
		db.SetMaxIdleConns(option.MaxIdleConn)
	}

	return &MysqlStorage{
		db: db,
	}, nil
}

func (ms *MysqlStorage) DB() *xorm.Engine {
	return ms.db
}
