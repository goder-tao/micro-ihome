package house

import (
	"github.com/gin-gonic/gin"
	"ihome/web/controller/user"
	"ihome/web/dao/mysql"
	house2 "ihome/web/logic/house"
	user2 "ihome/web/logic/user"
	"ihome/web/model/payload/house"
	"ihome/web/utils"
	"net/http"
	"strings"
)

// PostHouses 房东发布房源
func PostHouses(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	var housePl house.HousePayload
	if err := c.Bind(&housePl); err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "bind house: " + err.Error()
		return
	}
	houseId, err := house2.PostHouse(housePl, c)
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "bind house: " + err.Error()
		return
	}
	h["data"] = gin.H{"house_id": houseId}
}

// GetHouseInfo 获取房屋细节信息
func GetHouseInfo(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	hid := c.Param("id")
	hd, err := house2.GetHouseDetail(hid)
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "get house detail: " + err.Error()
		return
	}

	data := gin.H{
		"house":   hd,
		"user_id": hd.UserID,
	}
	h["data"] = data
}

// PostHousesImage 上传房屋图片
// 第一张上传的成为主图，并且注意持久化主图
func PostHousesImage(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)
	hid := c.Param("id")
	fh, err := c.FormFile("house_image")
	if err != nil {
		h["errno"] = utils.RECODE_USERERR
		h["errmsg"] = "form file: " + err.Error()
		return
	}
	f, err := fh.Open()
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "open file: " + err.Error()
		return
	}

	// 保存图片到fdfs
	buf := make([]byte, fh.Size)
	_, err = f.Read(buf)
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = " file read, " + err.Error()
		return
	}
	filePart := strings.Split(fh.Filename, ".")

	fileID, err := user2.SaveAvatarByte(buf, filePart[len(filePart)-1])
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "save house byte, " + err.Error()
		return
	}
	URL := user.URL_AVATAR_PREFIX + fileID
	h["data"] = gin.H{
		"url": URL,
	}

	// 可能是房屋的主图
	if err := mysql.TryStoreHouseURL(hid, URL); err != nil {
		h["errno"] = utils.RECODE_DATAERR
		h["errmsg"] = "store house index url, " + err.Error()
		return
	}
	if err := mysql.StoreHouseImg(hid, URL); err != nil {
		h["errno"] = utils.RECODE_DATAERR
		h["errmsg"] = "store house image url, " + err.Error()
		return
	}
}

// GetHouses 获取房源，用户搜索
func GetHouses(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	// aid(地区编号)、sd(开始时间)、ed(结束时间)、sk(排序类型)、page(页数)
	aid := c.Query("aid")
	sd := c.Query("sd")
	ed := c.Query("ed")
	sk := c.Query("sk")
	page := c.Query("p")

	//
	hs, tp, err := house2.SearchHouse(aid, sd, ed, sk, page)

	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "search house: " + err.Error()
		return
	}

	h["data"] = gin.H{
		"current_page": page,
		"houses":       hs,
		"total_page":   tp,
	}
}

// GetIndex 轮播房屋图片
func GetIndex(c *gin.Context) {

}
