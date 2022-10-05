package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/afocus/captcha"
	r "github.com/go-redis/redis/v8"
	"github.com/micro/micro/v3/service"
	capt "ihome/service/captcha/proto"
	"ihome/web/dao/redis"
	"log"
	"strings"
	"time"
)

func GetImageCd(uuid string) (*captcha.Image, error) {
	svc := service.New()
	cli := capt.NewCaptchaService("captcha", svc.Client())
	req := &capt.Request{}
	resp, err := cli.GetCaptcha(context.TODO(), req)

	if err != nil {
		return nil, err
	}
	var img captcha.Image
	if err := json.Unmarshal(resp.Data, &img); err != nil {
		return nil, err
	}

	// bind uuid and image code
	rdb, err := redis.Get()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	rdb.SetEX(ctx, uuid, resp.Code, time.Minute*5)

	return &img, nil
}

func CheckImageCd(uuid, inputImageCode string) error {
	rds, err := redis.Get()
	if err != nil {
		return err
	}
	ctx := context.Background()
	imageCode, err := rds.Get(ctx, uuid).Result()

	if err == r.Nil {
		return errors.New("image code expiration, please try again")
	} else if err != nil {
		return err
	} else if strings.ToLower(imageCode) != strings.ToLower(inputImageCode) {
		log.Println("imageCode: ", imageCode, "inputCode: ", inputImageCode)
		return errors.New("image code wrong")
	} else {
		return nil
	}
}
