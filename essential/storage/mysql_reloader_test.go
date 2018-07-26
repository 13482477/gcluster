package storage

import (
	"encoding/json"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type AdInfo struct {
	ID   int    `xorm:"'id'"`
	Name string `xorm:"'name'"`
}

type TestReloader struct {
	Data *AdInfo
}

func (r *TestReloader) Query(reloader *MysqlReloader) error {
	var data AdInfo
	_, err := reloader.DB.Table("ad").Where("id = 20000").Get(&data)

	if err != nil {
		return err
	}

	r.Data = &data
	return nil
}

func TestNew(t *testing.T) {
	db, err := xorm.NewEngine("mysql", "root:@(10.215.28.213:3306)/ad_fancy?charset=utf8")
	if err != nil {
		t.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("connected")

	reloader := new(TestReloader)

	NewMysqlReloader(db, reloader, 1)

	go func() {
		time.Sleep(time.Second * time.Duration(1))

		data, _ := json.Marshal(reloader)
		t.Log(string(data))
	}()

	t.Log("sleep")

	time.Sleep(time.Second * time.Duration(10))
}
