package dao

import (
	"github.com/LotteWong/giotto-gateway/po"
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type UserOperator struct{}

func NewUserOperator() *UserOperator {
	return &UserOperator{}
}

func (o *UserOperator) Find(c *gin.Context, tx *gorm.DB, req *po.Admin) (*po.Admin, error) {
	res := &po.Admin{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *UserOperator) Save(c *gin.Context, tx *gorm.DB, req *po.Admin) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(req).Error
	if err != nil {
		return err
	}
	return nil
}
