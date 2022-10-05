package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ihome/web/logic/user"
	"image/png"
	"net/http"
)

func GetImageCd(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	img, err := user.GetImageCd(uuid)

	if err != nil {
		fmt.Println("GetImage: ", err)
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	png.Encode(ctx.Writer, img)
	fmt.Println("uuid =", uuid)
}
