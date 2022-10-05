package user

import (
	"github.com/gin-gonic/gin"
	user2 "ihome/web/logic/user"
	"ihome/web/utils"
	"net/http"
)

func GetSession(ctx *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)

	name, err := user2.GetSessionName(ctx)
	if err != nil {
		h["errno"] = utils.RECODE_SESSIONERR
		h["errmsg"] = err.Error()
		ctx.JSON(http.StatusOK, h)
		return
	}
	h["data"] = map[string]string{"name": name}
	ctx.JSON(http.StatusOK, h)
}

// DeleteSession 用户退出删除session
func DeleteSession(c *gin.Context) {
	h := gin.H{}
	h["errno"] = utils.RECODE_OK
	h["errmsg"] = utils.RecodeText(utils.RECODE_OK)

	if err := user2.DeleteSession(c); err != nil {
		h["errno"] = utils.RECODE_SESSIONERR
		h["errmsg"] = "delete session, " + err.Error()
	}

	c.JSON(http.StatusOK, h)
}
