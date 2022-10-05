package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ihome/web/dao/mysql"
	"ihome/web/dao/redis"
	"ihome/web/logic/house"
	user2 "ihome/web/logic/user"
	"ihome/web/model/payload/user"
	"ihome/web/utils"
	"net/http"
	"strings"
)

const (
	FORM_AVATAR       = "avatar"
	URL_AVATAR_PREFIX = "http://192.168.31.142:8888/"
)

// PostRet 尝试注册一个用户
func PostRet(c *gin.Context) {
	pl := user.RegistryPayload{}
	ret := gin.H{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}
	if err := c.BindJSON(&pl); err != nil {
		ret["errno"] = utils.RECODE_DATAERR
		ret["errmsg"] = err.Error()
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	if err := user2.Register(pl.Mobile, pl.SmsCode, pl.Password); err != nil {
		ret["errno"] = utils.RECODE_USERERR
		ret["errmsg"] = err.Error()
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	c.JSON(http.StatusOK, ret)
}

func GetArea(c *gin.Context) {
	m := gin.H{}
	m["errno"] = utils.RECODE_OK
	m["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	code := http.StatusOK
	// 判断redis是否存在
	areas, err := redis.GetAreas()
	if err != nil {
		m["errno"] = utils.RECODE_DATAERR
		m["errmsg"] = "redis get area, " + err.Error()
		c.JSON(http.StatusBadRequest, m)
		return
	}
	if len(areas) != 0 {
		m["data"] = areas
	} else {
		// 不存在从数据库中取数据
		areas, err := mysql.GetArea()
		if err != nil {
			m["errno"] = utils.RECODE_DATAERR
			m["errmsg"] = "mysql get area, " + err.Error()
			code = http.StatusBadRequest
		}
		m["data"] = areas
		// 保存到redis
		if err := redis.SaveAreas(areas); err != nil {
			m["errno"] = utils.RECODE_DATAERR
			m["errmsg"] = "mysql save area, " + err.Error()
			code = http.StatusBadRequest
		}
	}

	c.JSON(code, m)
}

func PostLogin(c *gin.Context) {
	var pl user.LoginPayLoad
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	if err := c.Bind(&pl); err != nil {
		h["errno"] = utils.RECODE_DATAERR
		h["errmsg"] = "bind:" + err.Error()
		c.JSON(http.StatusBadRequest, h)
		return
	}

	// 1. 数据库校验用户正确性
	name, err := user2.CheckUser(pl.Mobile, pl.Password)
	if err != nil {
		h["errno"] = utils.RECODE_DATAERR
		h["errmsg"] = "checkuser:" + err.Error()
		c.JSON(http.StatusBadRequest, h)
		return
	}
	// 2. 存入session
	if err := user2.SaveSessionName(c, name); err != nil {
		h["errno"] = utils.RECODE_SESSIONERR
		h["errmsg"] = "session save name:" + err.Error()
		c.JSON(http.StatusBadRequest, h)
		return
	}

	c.JSON(http.StatusOK, h)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	// 从session获取当前用户名
	userName, err := user2.GetSessionName(c)
	if err != nil {
		h["errno"] = utils.RECODE_SESSIONERR
		h["errmsg"] = "get session, " + err.Error()
		return
	}

	// 到mysql中查询用户信息
	u, err := user2.GetUserInfo(userName)
	if err != nil {
		h["errno"] = utils.RECODE_DATAERR
		h["errmsg"] = "get userinfo from db, " + err.Error()
		return
	}
	h["data"] = u
}

func PutUserName(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	var name struct {
		Name string `json:"name"`
	}

	if err := c.Bind(&name); err != nil {
		h["errno"] = utils.RECODE_DATAERR
		h["errmsg"] = "bind, " + err.Error()
		return
	} else if name.Name == "" {
		h["errno"] = utils.RECODE_USERERR
		h["errmsg"] = "forbid empty name"
		return
	}

	if err := user2.PutUserName(c, name.Name); err != nil {
		h["errno"] = utils.RECODE_DATAERR
		h["errmsg"] = "put user name, " + err.Error()
		return
	}

	h["data"] = name
}

// PostAvatar 上传用户头像
func PostAvatar(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	fh, err := c.FormFile(FORM_AVATAR)
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "form file, " + err.Error()
		return
	}

	// 1. 存储头像
	f, err := fh.Open()
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = " fh open file, " + err.Error()
		return
	}
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
		h["errmsg"] = "save avatar byte, " + err.Error()
		return
	}

	avatarURL := URL_AVATAR_PREFIX + fileID
	// 2. 持久化头像的file_id
	if err := user2.SaveAvatarURL(c, avatarURL); err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "save avatar url, " + err.Error()
		return
	}
	// err = c.SaveUploadedFile(fh, "/home/tao/Data/Software/project/go/project/IHome/web/test/web/img.png")
	h["data"] = gin.H{"avatar_url": avatarURL}
}

// PutUserAuth 用户身份实名认证
func PutUserAuth(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	var realName struct {
		Id_card   string `json:"id_card"`
		Real_name string `json:"real_name"`
	}
	if err := c.Bind(&realName); err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "bind data: " + err.Error()
		return
	}
	if err := user2.AuthUser(realName.Real_name, realName.Id_card, c); err != nil {
		h["errno"] = utils.RECODE_USERERR
		h["errmsg"] = "real name auth fail: " + err.Error()
	}
}

func GetUserAuth(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	u, err := user2.GetUserAuth(c)
	if err != nil {
		h["errno"] = utils.RECODE_OK
		h["errmsg"] = "get user auth: " + err.Error()
		return
	}
	h["data"] = u
}

// GetUserHouses 获取房东发布的房源信息
func GetUserHouses(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	defer c.JSON(http.StatusOK, h)

	oname, err := user2.GetSessionName(c)
	if err != nil {
		h["errno"] = utils.RECODE_SESSIONERR
		h["errmsg"] = "get session name: " + err.Error()
		return
	}

	ohs, err := house.GetOwnerHouses(oname)
	if err != nil {
		h["errno"] = utils.RECODE_UNKNOWERR
		h["errmsg"] = "get user house: " + err.Error()
		return
	}

	h["data"] = ohs
}

// GetUserOrder 根据角色获取订单状态（无差别）
func GetUserOrder(c *gin.Context) {
	role := c.Query("role")
	fmt.Println(role)
}
