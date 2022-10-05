package order

import (
	"ihome/web/dao/mysql"
	"ihome/web/model/payload/order"
)

func PostOrders(opl order.OrderPayload) (int, error) {
	return mysql.PostOrders(opl)
}

func PutOrder(oid uint, status string) error {
	return mysql.PutOrder(oid, status)
}

func PutComment(oid uint, comment string) error {
	return mysql.PutComment(oid, comment)
}
