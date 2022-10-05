package house

// HousePayload 发布房源信息
type HousePayload struct {
	Title     string   `json:"title"`
	Price     string   `json:"price"`
	AreaID    string   `json:"area_id"`
	Address   string   `json:"address"`
	RoomCount string   `json:"room_count"`
	Acreage   string   `json:"acreage"`
	Unit      string   `json:"unit"`
	Capacity  string   `json:"capacity"`
	Beds      string   `json:"beds"`
	Deposit   string   `json:"deposit"`
	MinDays   string   `json:"min_days"`
	MaxDays   string   `json:"max_days"`
	Facility  []string `json:"facility"`
}

// HouseDetail 房屋详细信息
type HouseDetail struct {
	HousePayload
	// Comment
	Hid        string     `json:"hid"`
	ImageURLs  []string   `json:"img_urls"`
	UserID     string     `json:"user_id"`
	UserName   string     `json:"user_name"`
	UserAvatar string     `json:"user_avatar"`
	Comments   []CommentT `json:"comments"`
}
type CommentT struct {
	Comment  string `json:"comment"`
	CTime    string `json:"ctime"`
	UserName string `json:"user_name"`
}

// HouseSourceInfo 房源信息，用户搜索房屋时用
type HouseSourceInfo struct {
	Address    string `json:"address"`
	AreaName   string `json:"area_name"`
	Ctime      string `json:"ctime"`
	HouseId    int    `json:"house_id"`
	ImgURL     string `json:"img_url"`
	OrderCount int    `json:"order_count"`
	Price      int    `json:"price"`
	RoomCount  int    `json:"room_count"`
	Title      string `json:"title"`
	UserAvatar string `json:"user_avatar"`
}

// OwnerHouse 房东获取房源信息
type OwnerHouse struct {
	Address    string `json:"address"`
	AreaName   string `json:"area_name"`
	CTime      string `json:"ctime"`
	HouseID    uint   `json:"house_id"`
	ImgURL     string `json:"img_url"`
	OrderCount int    `json:"order_count"`
	Price      int    `json:"price"`
	RoomCount  int    `json:"room_count"`
	Title      string `json:"title"`
	UserAvatar string `json:"user_avatar"`
}
