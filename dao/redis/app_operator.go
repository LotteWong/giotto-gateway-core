package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/utils"
	"github.com/garyburd/redigo/redis"
)

type AppOperator struct{}

func NewAppOperator() *AppOperator {
	return &AppOperator{}
}

func (o *AppOperator) GetApp(appId int64) (*po.App, error) {
	data, err := redis.Bytes(RedisConfDo("get", fmt.Sprintf("%s_%d", "app", appId)))
	if err != nil {
		return nil, err
	}
	res := &po.App{}
	json.Unmarshal(data, res)
	return res, nil
}

func (o *AppOperator) SetApp(appId int64, obj *po.App) error {
	data := utils.Obj2Json(obj)
	_, err := RedisConfDo("set", fmt.Sprintf("%s_%d", "app", appId), data)
	if err != nil {
		return err
	}
	return nil
}

func (o *AppOperator) DelApp(appId int64) error {
	_, err := RedisConfDo("del", fmt.Sprintf("%s_%d", "app", appId))
	if err != nil {
		return err
	}
	return nil
}

func (o *AppOperator) ListApps() ([]*po.App, error) {
	var apps []*po.App

	keys, err := redis.ByteSlices(RedisConfDo("keys", "app_*"))
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		appId, err := strconv.Atoi(strings.ReplaceAll(string(key), "app_", ""))
		if err != nil {
			return nil, err
		}

		app, err := o.GetApp(int64(appId))
		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}

	return apps, nil
}
