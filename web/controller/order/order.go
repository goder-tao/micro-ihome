package order

import (
	"github.com/gin-gonic/gin"
	order2 "ihome/web/logic/order"
	"ihome/web/model/payload/order"
	"ihome/web/utils"
	"net/http"
	"strconv"
)

// PostOrders 提交订单
func PostOrders(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)
	var opl order.OrderPayload

	if err := c.Bind(&opl); err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "bind: " + err.Error()
		return
	}
	order_id, err := order2.PostOrders(opl)
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "post order: " + err.Error()
		return
	}
	h["data"] = order_id
}

func PutOrders(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)
	id := c.Param("id")
	oid, _ := strconv.Atoi(id)
	var action order.OrderStatusPayload
	if err := c.Bind(&action); err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "bind: " + err.Error()
		return
	}
	if err := order2.PutOrder(uint(oid), action.Action); err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "put order: " + err.Error()
	}
}

func PutComment(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)
	var cpl order.OrderCommentPayload
	if err := c.Bind(&cpl); err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "bind: " + err.Error()
		return
	}
	if err := order2.PutComment(cpl.OrderId, cpl.Comment); err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "put comment: " + err.Error()
		return
	}
}
