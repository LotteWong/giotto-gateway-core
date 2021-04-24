package mysql

import (
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceOperator struct{}

func NewServiceOperator() *ServiceOperator {
	return &ServiceOperator{}
}

func (o *ServiceOperator) Find(ctx *gin.Context, tx *gorm.DB, req *po.ServiceInfo) (*po.ServiceInfo, error) {
	res := &po.ServiceInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *ServiceOperator) Save(ctx *gin.Context, tx *gorm.DB, req *po.ServiceInfo) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Save(req).Error; err != nil {
		return err
	}
	return nil
}

func (o *ServiceOperator) FuzzySearchAndPage(ctx *gin.Context, tx *gorm.DB, keyword string, pageIndex, pageSize int) (int64, []po.ServiceInfo, error) {
	var total int64
	var items []po.ServiceInfo

	query := tx.SetCtx(utils.GetGinTraceContext(ctx)).Table((&po.ServiceInfo{}).TableName()).Where("is_delete=0")

	if keyword != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if pageSize != 0 {
		query = query.Limit(pageSize)
	}
	if pageIndex != 0 {
		query = query.Offset((pageIndex - 1) * pageSize)
	}

	if err := query.Order("id desc").Find(&items).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, nil, err
	}
	query.Limit(-1).Offset(-1).Count(&total)

	return total, items, nil
}

func (o *ServiceOperator) GroupByServiceType(ctx *gin.Context, tx *gorm.DB) ([]po.ServicePercentage, error) {
	var groups []po.ServicePercentage

	query := tx.SetCtx(utils.GetGinTraceContext(ctx)).Table((&po.ServiceInfo{}).TableName()).Where("is_delete=0")
	if err := query.Select("service_type, count(*) as service_count").Group("service_type").Scan(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}
