package mysql

import (
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type LoadBalanceOperator struct{}

func NewLoadBalanceOperator() *LoadBalanceOperator {
	return &LoadBalanceOperator{}
}

func (o *LoadBalanceOperator) Find(ctx *gin.Context, tx *gorm.DB, req *po.LoadBalance) (*po.LoadBalance, error) {
	res := &po.LoadBalance{}
	err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *LoadBalanceOperator) Save(ctx *gin.Context, tx *gorm.DB, req *po.LoadBalance) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Save(req).Error; err != nil {
		return err
	}
	return nil
}

func (o *LoadBalanceOperator) ListByServiceId(ctx *gin.Context, tx *gorm.DB, serviceId int64) ([]po.LoadBalance, int64, error) {
	var loadBalanceList []po.LoadBalance
	var count int64
	var err error
	query := tx.SetCtx(utils.GetGinTraceContext(ctx)).Table((&po.LoadBalance{}).TableName()).Where("service_id=?", serviceId)

	err = query.Order("id desc").Find(&loadBalanceList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return loadBalanceList, count, nil
}
