package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"ihome/web/model"
	"strings"
	"time"
)

func SaveSms(phone, smsCode string) error {
	// 存入smscode存入redis
	rds, err := Get()
	if err != nil {
		return err
	}
	rds.SetEX(context.Background(), phone, smsCode, time.Minute*5)
	return nil
}

func CheckSms(phone, inputSms string) error {
	rdb, err := Get()
	if err != nil {
		return err
	}
	smsCode, err := rdb.Get(context.Background(), phone).Result()

	if err == redis.Nil {
		return errors.New("sms code expiration, please send again")
	} else if err != nil {
		return err
	} else if strings.ToLower(inputSms) != strings.ToLower(smsCode) {
		return errors.New("sms code wrong!")
	}
	return nil
}

func GetAreas() ([]model.Area, error) {
	var areas []model.Area
	rdb, err := Get()
	if err != nil {
		return nil, err
	}
	areaNames, err := rdb.ZRange(context.Background(), KEY_AREA, 0, rdb.ZCard(context.Background(), KEY_AREA).Val()+1).Result()
	if err != nil {
		return nil, err
	}
	for i, name := range areaNames {
		areas = append(areas, model.Area{
			Id:   i + 1,
			Name: name,
		})
	}
	return areas, nil
}

func SaveAreas(areas []model.Area) error {
	rdb, err := Get()
	if err != nil {
		return err
	}
	var zs []*redis.Z
	for _, area := range areas {
		zs = append(zs, &redis.Z{
			Score:  float64(area.Id),
			Member: area.Name,
		})
	}
	rdb.ZAdd(context.Background(), KEY_AREA, zs...)
	return nil
}

func GetUserAvatarURL(userID string) (string, error) {
	rdb, err := Get()
	if err != nil {
		return "", err
	}
	avatar, err := rdb.Get(context.Background(), KEY_PREFIX_USERAVATAR+userID).Result()
	if err == redis.Nil {
		return " ", nil
	}
	return avatar, nil
}

func SaveUserAvatarURL(userID, avatarURL string) error {
	// 临时存储一分钟
	rdb, err := Get()
	if err != nil {
		return err
	}
	return rdb.SetNX(context.Background(), KEY_PREFIX_USERAVATAR+userID, avatarURL, EXPIRATION_USERAVATAR).Err()
}
