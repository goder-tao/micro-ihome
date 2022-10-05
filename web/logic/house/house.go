package house

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ihome/web/dao/mysql"
	"ihome/web/dao/redis"
	"ihome/web/logic/user"
	"ihome/web/model"
	"ihome/web/model/payload/house"
	"strconv"
)

const (
	// search kind
	NEW     = "new"
	PRICE   = "price"
	ACREAGE = "acreage"
	ORDER   = "order_count"

	//
	PER_PAGE = 10
)

var searchKind = map[string]string{
	NEW:     "created_at",
	PRICE:   "price",
	ACREAGE: "acreage",
	ORDER:   "order_count",
}

func PostHouse(pl house.HousePayload, c *gin.Context) (int, error) {
	// pl转换成model
	name, err := user.GetSessionName(c)
	if err != nil {
		return 0, err
	}
	var h model.House
	u, err := mysql.GetUserInfoByName(name)
	if err != nil {
		return 0, err
	}
	h.UserId = uint(u.ID)
	aid, _ := strconv.Atoi(pl.AreaID)
	h.AreaId = uint(aid)
	h.Title = pl.Title
	h.Address = pl.Address
	h.Room_count, _ = strconv.Atoi(pl.RoomCount)
	h.Acreage, _ = strconv.Atoi(pl.Acreage)
	h.Price, _ = strconv.Atoi(pl.Price)
	h.Unit = pl.Unit
	h.Capacity, _ = strconv.Atoi(pl.Capacity)
	h.Beds = pl.Beds
	h.Deposit, _ = strconv.Atoi(pl.Deposit)
	h.Min_days, _ = strconv.Atoi(pl.MinDays)
	h.Max_days, _ = strconv.Atoi(pl.MaxDays)

	// 尝试从redis取基础设施
	fids := pl.Facility
	fcs, err := redis.GetFacilities(fids)
	if err != nil {
		return 0, err
	}
	// 基础设施取不到从mysql取并且存入redis中
	if len(fids) != 0 && len(fcs) == 0 {
		fcsa, err := mysql.GetFacilities(nil)
		if err != nil {
			return 0, err
		}
		err = redis.SaveFacilities(fcsa)
		if err != nil {
			return 0, err
		}
		fcs, err = redis.GetFacilities(fids)
		if err != nil {
			return 0, err
		}
	}
	h.Facilities = fcs
	// 持久化
	return mysql.PostHouse(h)
}

// SearchHouse 租客查找房源
func SearchHouse(aid, sd, ed, sk, page string) ([]house.HouseSourceInfo, int, error) {
	if _, ok := searchKind[sk]; !ok {
		return nil, 0, errors.New("wrong search kind")
	}
	// 数据库中查找到初步数据
	p, _ := strconv.Atoi(page)
	ret, cnt, err := mysql.SearchHouse(aid, sd, ed, searchKind[sk], (p-1)*PER_PAGE, PER_PAGE)
	if err != nil {
		return nil, 0, err
	}

	// 处理house的其他数据
	for i := 0; i < len(ret); i++ {
		areaName, err := redis.GetAreaName(ret[i].AreaName)
		if err != nil {
			return nil, 0, err
		}
		// redis中替换AreaName
		ret[i].AreaName = areaName
		// avatar_url
		avatar, err := redis.GetUserAvatarURL(ret[i].UserAvatar)
		if err != nil {
			return nil, 0, err
		}
		// 没有avatar尝试向mysql取
		if avatar == " " {
			avatar, err = mysql.GetUserAvatarURL(ret[i].UserAvatar)
			if err != nil {
				return nil, 0, err
			}
			err = redis.SaveUserAvatarURL(ret[i].UserAvatar, avatar)
		}
		ret[i].UserAvatar = avatar
	}

	tp := cnt / PER_PAGE
	if cnt%PER_PAGE != 0 {
		tp += 1
	}
	return ret, tp, nil
}

// GetHouseDetail 获取房屋细节信息，填充HouseDetail
func GetHouseDetail(hid string) (house.HouseDetail, error) {
	var hd house.HouseDetail
	baseInfo, err := mysql.GetHouseInfo(hid)
	if err != nil {
		return house.HouseDetail{}, err
	}
	u, err := mysql.GetUserInfoById(strconv.Itoa(int(baseInfo.UserId)))
	if err != nil {
		return house.HouseDetail{}, err
	}
	// house_img表
	urls, err := mysql.GetHouseImages(hid)
	if err != nil {
		return house.HouseDetail{}, err
	}
	// OrderHouse表
	comments, err := mysql.GetHouseComment(hid)
	if err != nil {
		return house.HouseDetail{}, err
	}
	// facility
	fcs, err := mysql.GetHouseFacility(hid)

	// 信息拼接
	hd.Hid = hid
	hd.UserID = strconv.Itoa(u.ID)
	hd.UserName = u.Name
	hd.UserAvatar = u.Avatar_url
	hd.Comments = comments
	hd.ImageURLs = urls
	hd.Title = baseInfo.Title
	hd.Price = strconv.Itoa(baseInfo.Price)
	hd.Address = baseInfo.Address
	hd.RoomCount = strconv.Itoa(baseInfo.Room_count)
	hd.Acreage = strconv.Itoa(baseInfo.Acreage)
	hd.Unit = baseInfo.Unit
	hd.Capacity = strconv.Itoa(baseInfo.Capacity)
	hd.Beds = baseInfo.Beds
	hd.Deposit = strconv.Itoa(baseInfo.Deposit)
	hd.MinDays = strconv.Itoa(baseInfo.Min_days)
	hd.MaxDays = strconv.Itoa(baseInfo.Max_days)
	hd.Facility = fcs
	return hd, nil
}

func GetOwnerHouses(oname string) ([]house.OwnerHouse, error) {
	u, err := mysql.GetUserInfoByName(oname)
	if err != nil {
		return nil, err
	}
	ohs, err := mysql.GetOwnerHouses(u.ID)
	if err != nil {
		return nil, err
	}
	for _, oh := range ohs {
		oh.UserAvatar = u.Avatar_url
	}
	return ohs, nil
}
