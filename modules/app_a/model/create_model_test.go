package model

import (
	"testing"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"time"
)

func TestCreateModel(t *testing.T) {
	db, err := gorm.Open("mysql", "root:password@(localhost:3306)/mcloud?parseTime=true&loc=Local")
	if err != nil {
		log.Panicf("Fatal error mysql connection failed: %v", err)
	}
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)           // 设置连接池大小
	db.DB().SetMaxOpenConns(100)          // 设置最大连接数
	db.DB().SetConnMaxLifetime(time.Hour) // 设置连接可被复用的最大时长

	if err := db.AutoMigrate(&Ecs{}).Error; err != nil {
		log.Error(err)
	}

}
