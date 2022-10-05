package mysql

import (
	"ihome/web/model"
	"ihome/web/model/payload/order"
	"time"
)

func PostOrders(opl order.OrderPayload) (int, error) {
	db, err := Get(IHOME)
	if err != nil {
		return 0, err
	}
	sd, err := time.Parse("2006-01-02 15:04:05", opl.StartDate)
	ed, err := time.Parse("2006-01-02 15:04:05", opl.EndDate)
	if err != nil {
		return 0, err
	}
	oh := model.OrderHouse{HouseId: opl.HouseID, Begin_date: sd, End_date: ed}
	ret := db.Create(&oh)
	if ret.Error != nil {
		return 0, nil
	}
	return int(oh.ID), nil
}

func PutOrder(oid uint, status string) error {
	db, err := Get(IHOME)
	if err != nil {
		return err
	}
	return db.Where("id = ?", oid).Update("status", status).Error
}

func PutComment(oid uint, comment string) error {
	db, err := Get(IHOME)
	if err != nil {
		return err
	}
	return db.Where("id = ?", oid).Update("comment", comment).Error
}
