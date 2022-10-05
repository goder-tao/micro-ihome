package handler

import (
	"context"
	"encoding/json"
	capt "github.com/afocus/captcha"
	_ "github.com/micro/micro/v3/service/logger"
	captcha "ihome/service/captcha/proto"
	"image/color"
)

type Captcha struct{}

func (e *Captcha) GetCaptcha(ctx context.Context, request *captcha.Request, response *captcha.Response) error {
	c := capt.New()

	if err := c.SetFont("/home/tao/Data/Software/project/go/project/IHome/web/config/comic.ttf"); err != nil {
		return err
	}

	c.SetSize(128, 64)

	c.SetDisturbance(capt.MEDIUM)

	c.SetFrontColor(color.RGBA{255, 255, 255, 255})
	c.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	img,code := c.Create(4,capt.ALL)

	data, _ := json.Marshal(img)

	response.Data = data
	response.Code = code
	return nil
}

// Return a new handler
func New() *Captcha {
	return &Captcha{}
}