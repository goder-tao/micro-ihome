package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ihome/web/model"
	"log"
	"sync"
)

const (
	user     = "root"
	password = "password"
	address  = "127.0.0.1:3306"
	IHOME    = "ihome"
)

var dbs map[string]*gorm.DB = make(map[string]*gorm.DB)
var mus map[string]*sync.Mutex = make(map[string]*sync.Mutex)
var mu sync.Mutex

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, address, IHOME)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("mysql init fail: ", err)
	}
	dbs[IHOME] = db
	// new(model.House), new(model.User), new(model.Area), new(model.Facility), new(model.HouseImage), new(model.OrderHouse)
	if err := db.AutoMigrate(new(model.House), new(model.User), new(model.Area), new(model.Facility), new(model.HouseImage), new(model.OrderHouse)); err != nil {
		log.Fatalln("mysql aotumigrate fail: ", err)
	}
}

func Get(dbName string) (*gorm.DB, error) {
	if dbs[dbName] == nil {
		if _, ok := mus[dbName]; !ok {
			mu.Lock()
			mus[dbName] = &sync.Mutex{}
			mu.Unlock()
		}
		mus[dbName].Lock()
		defer mus[dbName].Unlock()

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, address, dbName)
		// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		dbs[dbName] = db
	}
	return dbs[dbName], nil
}
