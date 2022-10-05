package order

// OrderPayload 提交订单payload
type OrderPayload struct {
	HouseID   uint   `json:"house_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type OrderStatusPayload struct {
	Action string `json:"action"`
}

type OrderCommentPayload struct {
	OrderId uint   `json:"order_id"`
	Comment string `json:"comment"`
}
