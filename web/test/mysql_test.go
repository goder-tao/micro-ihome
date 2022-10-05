package test

import (
	"ihome/web/dao/mysql"
	"testing"
)

type Test struct {
	Id int
	T  int
}

func TestMysql(t *testing.T) {
	//db, err := mysql.Get("ihome")
	//if err != nil {
	//	t.Error("db get: ", err)
	//}
	//
	//err = db.AutoMigrate(new(Test))
	//if err != nil {
	//	t.Error("create table: ", err)
	//}
	//db.Create(&Test{Id: 1, T: 1})
	//var ts []Test
	//db.Find(&ts)
	//for i := 0; i < len(ts); i++ {
	//	fmt.Println(ts[i].Id, ts[i].T)
	//}
}

func TestInit(t *testing.T) {
	mysql.Init()
}
