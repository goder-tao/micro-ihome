package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"ihome/web/model"
	"strconv"
)

func GetFacilities(ids []string) ([]*model.Facility, error) {
	// hash结构
	cli, err := Get()
	if err != nil {
		return nil, err
	}
	var ret []*model.Facility
	for _, id := range ids {
		fc := &model.Facility{}
		fc.Name = cli.HGet(context.Background(), KEY_FACILITY, id).Val()
		fc.Id, _ = strconv.Atoi(id)
		ret = append(ret, fc)
	}
	return ret, nil
}

// SaveFacilities 缓存基础设施
func SaveFacilities(fcs []*model.Facility) error {
	cli, err := Get()
	if err != nil {
		return err
	}

	for _, fc := range fcs {
		cli.HSet(context.Background(), KEY_FACILITY, fc.Id, fc.Name)
	}
	return nil
}

func GetAreaName(areaId string) (string, error) {
	cli, err := Get()
	if err != nil {
		return "", err
	}
	vs, err := cli.ZRangeByScore(context.Background(), KEY_AREA, &redis.ZRangeBy{Min: areaId, Max: areaId}).Result()
	if err != nil {
		return "", err
	}
	return vs[0], nil
}
