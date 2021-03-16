package dao

import (
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type AppOperator struct{}

func NewAppOperator() *AppOperator {
	return &AppOperator{}
}

func (o *AppOperator) Find(ctx *gin.Context, tx *gorm.DB, req *po.App) (*po.App, error) {
	res := &po.App{}
	err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *AppOperator) Save(ctx *gin.Context, tx *gorm.DB, req *po.App) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Save(req).Error; err != nil {
		return err
	}
	return nil
}

func (o *AppOperator) FuzzySearchAndPage(ctx *gin.Context, tx *gorm.DB, keyword string, pageIndex, pageSize int) (int64, []po.App, error) {
	var total int64
	var items []po.App

	query := tx.SetCtx(utils.GetGinTraceContext(ctx)).Table((&po.App{}).TableName()).Where("is_delete=0")

	if keyword != "" {
		query = query.Where("(app_id like ? or app_name like ?)", "%"+keyword+"%", "%"+keyword+"%")
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
