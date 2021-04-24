package mysql

import (
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type AccessControlOperator struct{}

func NewAccessControlOperator() *AccessControlOperator {
	return &AccessControlOperator{}
}

func (o *AccessControlOperator) Find(ctx *gin.Context, tx *gorm.DB, req *po.AccessControl) (*po.AccessControl, error) {
	res := &po.AccessControl{}
	err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *AccessControlOperator) Save(ctx *gin.Context, tx *gorm.DB, req *po.AccessControl) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Save(req).Error; err != nil {
		return err
	}
	return nil
}

func (o *AccessControlOperator) ListByServiceId(ctx *gin.Context, tx *gorm.DB, serviceId int64) ([]po.AccessControl, int64, error) {
	var accessControlList []po.AccessControl
	var count int64
	var err error
	query := tx.SetCtx(utils.GetGinTraceContext(ctx)).Table((&po.AccessControl{}).TableName()).Where("service_id=?", serviceId)

	err = query.Order("id desc").Find(&accessControlList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return accessControlList, count, nil
}
