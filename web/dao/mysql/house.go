package mysql

import (
	"ihome/web/dao/redis"
	"ihome/web/model"
	"ihome/web/model/payload/house"
	"strconv"
)

func PostHouse(h model.House) (int, error) {
	db, err := Get(IHOME)
	if err != nil {
		return 0, err
	}
	err = db.Create(&h).Error
	if err != nil {
		return 0, err
	}
	return int(h.ID), nil
}

// GetFacilities 获取基础设施的信息
// fids 支持返回特定的基础设施，空值默认返回所有
func GetFacilities(fids []string) ([]*model.Facility, error) {
	var ret []*model.Facility
	db, err := Get(IHOME)
	if err != nil {
		return nil, err
	}
	if len(fids) == 0 {
		// get all
		rows, err := db.Model(model.Facility{}).Rows()
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			fc := &model.Facility{}
			err := rows.Scan(fc)
			if err != nil {
				return nil, err
			}
			ret = append(ret, fc)
		}
	} else {

	}
	return ret, nil
}

func TryStoreHouseURL(hid, URL string) error {
	db, err := Get(IHOME)
	if err != nil {
		return err
	}
	return db.Model(model.House{}).Where("id = ? AND index_image_url = ?", hid, "").Update("index_image_url", URL).Error
}

func SearchHouse(aid, sd, ed, sk string, offset, limit int) ([]house.HouseSourceInfo, int, error) {
	db, err := Get(IHOME)
	if err != nil {
		return nil, 0, err
	}
	rows, err := db.Model(model.House{}).Offset(offset).Limit(limit).Order(sk).Where("area_id = ?", aid).Rows()
	var cnt int64
	db.Model(model.House{}).Offset(-1).Limit(-1).Count(&cnt)
	if err != nil {
		return nil, 0, err
	}
	var ret []house.HouseSourceInfo

	for rows.Next() {
		var h model.House
		var hsi house.HouseSourceInfo
		if err := db.ScanRows(rows, &h); err != nil {
			return nil, 0, err
		}

		hsi.HouseId = int(h.ID)
		hsi.Price = h.Price
		hsi.RoomCount = h.Room_count
		hsi.Address = h.Address
		hsi.Title = h.Title
		hsi.AreaName = strconv.Itoa(int(h.AreaId))
		hsi.Ctime = h.CreatedAt.String()
		hsi.OrderCount = h.Order_count
		hsi.UserAvatar = strconv.Itoa(int(h.UserId))
		hsi.ImgURL = h.Index_image_url

		ret = append(ret, hsi)
	}
	return ret, int(cnt), nil
}

func StoreHouseImg(hids, URL string) error {
	db, err := Get(IHOME)
	if err != nil {
		return err
	}
	hid, _ := strconv.Atoi(hids)
	return db.Create(&model.HouseImage{HouseId: uint(hid), Url: URL}).Error
}

// GetHouseInfo 获取房屋的基本信息（model.House）
func GetHouseInfo(hid string) (model.House, error) {
	db, err := Get(IHOME)
	if err != nil {
		return model.House{}, err
	}
	var h model.House

	err = db.Where("id = ?", hid).Find(&h).Error
	if err != nil {
		return model.House{}, err
	}
	return h, nil
}

// GetHouseImages 获取房屋所有的图片url(model.HouseImage)
func GetHouseImages(hid string) ([]string, error) {
	db, err := Get(IHOME)
	if err != nil {
		return nil, err
	}
	var ret []string
	var store model.HouseImage
	rows, err := db.Model(model.HouseImage{}).Where("house_id = ?", hid).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = db.ScanRows(rows, &store)
		if err != nil {
			return nil, err
		}
		ret = append(ret, store.Url)
	}
	return ret, nil
}

// GetHouseComment 获取房屋的评论信息(model.OrderHouse)
func GetHouseComment(hid string) ([]house.CommentT, error) {
	db, err := Get(IHOME)
	if err != nil {
		return nil, err
	}
	var ret []house.CommentT
	rows, err := db.Model(model.OrderHouse{}).Select("order_houses.comment, order_houses.created_at, users.name").Joins("left join users on order_houses.user_id = users.id").Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := house.CommentT{}
		if err := rows.Scan(&s.Comment, &s.CTime, &s.UserName); err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}

	return ret, nil
}

// GetHouseFacility 获取房屋的基础设施id
func GetHouseFacility(hid string) ([]string, error) {
	db, err := Get(IHOME)
	if err != nil {
		return nil, err
	}
	rows, err := db.Table("houses_facilities").Select("facility_id").Where("house_id = ?", hid).Rows()
	if err != nil {
		return nil, err
	}
	var ret []string
	var s string
	for rows.Next() {
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func GetOwnerHouses(uid int) ([]house.OwnerHouse, error) {
	// house id: *cnt
	m := make(map[uint]*int)
	db, err := Get(IHOME)
	if err != nil {
		return nil, err
	}
	// houses left join order_houses
	rows, err := db.Model(model.House{}).Joins("left join order_houses on houses.id=order_houses.house_id").Where("houses.user_id = ?", uid).Rows()
	if err != nil {
		return nil, err
	}

	var ret []house.OwnerHouse
	for rows.Next() {
		h := model.House{}
		if err := db.ScanRows(rows, &h); err != nil {
			return nil, err
		}
		if _, ok := m[h.ID]; !ok {
			m[h.ID] = new(int)
			aname, err := redis.GetAreaName(strconv.Itoa(int(h.AreaId)))
			if err != nil {
				return nil, err
			}
			oh := house.OwnerHouse{
				Address:    h.Address,
				AreaName:   aname,
				CTime:      h.CreatedAt.String(),
				HouseID:    h.ID,
				ImgURL:     h.Index_image_url,
				OrderCount: *m[h.ID],
				Price:      h.Price,
				RoomCount:  h.Room_count,
				Title:      h.Title,
				UserAvatar: "",
			}
			ret = append(ret, oh)
		}
		*m[h.ID] += 1
	}
	return ret, nil
}
