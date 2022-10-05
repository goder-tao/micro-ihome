package user

import (
	"github.com/gin-gonic/gin"
	"ihome/web/logic/user"
	"ihome/web/utils"
	"net/http"
)

func GetSmscd(ctx *gin.Context) {
	phone := ctx.Param("mobile")
	inputImgCode := ctx.Query("text")
	uuid := ctx.Query("id")

	ret := gin.H{
		"errno":  "",
		"errmsg": "",
	}

	// 1. 图片验证码校验
	err := user.CheckImageCd(uuid, inputImgCode)
	if err == nil {
		// 2. 短信验证码发送（模拟）
		err := user.SendSms(phone)
		if err != nil {
			ret["errno"] = utils.RECODE_SMSERR
			ret["errmsg"] = err.Error()
		} else {
			ret["errno"] = utils.RECODE_OK
			ret["errmsg"] = utils.RecodeText(utils.RECODE_OK)
		}
	} else {
		// 图片验证失败
		ret["errno"] = utils.RECODE_DATAERR
		ret["errmsg"] = err.Error()
	}
	ctx.JSONP(http.StatusOK, ret)
}
