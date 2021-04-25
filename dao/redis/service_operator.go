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

type ServiceOperator struct{}

func NewServiceOperator() *ServiceOperator {
	return &ServiceOperator{}
}

func (o *ServiceOperator) GetService(serviceId int64) (*po.ServiceDetail, error) {
	data, err := redis.Bytes(RedisConfDo("get", fmt.Sprintf("%s_%d", "service", serviceId)))
	if err != nil {
		return nil, err
	}
	res := &po.ServiceDetail{}
	json.Unmarshal(data, res)
	return res, nil
}

func (o *ServiceOperator) SetService(serviceId int64, obj *po.ServiceDetail) error {
	data := utils.Obj2Json(obj)
	_, err := RedisConfDo("set", fmt.Sprintf("%s_%d", "service", serviceId), data)
	if err != nil {
		return err
	}
	return nil
}

func (o *ServiceOperator) DelService(serviceId int64) error {
	_, err := RedisConfDo("del", fmt.Sprintf("%s_%d", "service", serviceId))
	if err != nil {
		return err
	}
	return nil
}

func (o *ServiceOperator) ListServices() ([]*po.ServiceDetail, error) {
	var services []*po.ServiceDetail

	keys, err := redis.ByteSlices(RedisConfDo("keys", "service_*"))
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		serviceId, err := strconv.Atoi(strings.ReplaceAll(string(key), "service_", ""))
		if err != nil {
			return nil, err
		}

		service, err := o.GetService(int64(serviceId))
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	return services, nil
}
