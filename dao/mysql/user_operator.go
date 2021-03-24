package mysql

import (
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type UserOperator struct{}

func NewUserOperator() *UserOperator {
	return &UserOperator{}
}

func (o *UserOperator) Find(ctx *gin.Context, tx *gorm.DB, req *po.Admin) (*po.Admin, error) {
	res := &po.Admin{}
	err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *UserOperator) Save(ctx *gin.Context, tx *gorm.DB, req *po.Admin) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Save(req).Error; err != nil {
		return err
	}
	return nil
}
