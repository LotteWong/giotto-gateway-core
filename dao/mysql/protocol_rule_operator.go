package mysql

import (
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ProtocolRuleOperator struct{}

func NewProtocolRuleOperator() *ProtocolRuleOperator {
	return &ProtocolRuleOperator{}
}

func (o *ProtocolRuleOperator) FindHttpRule(ctx *gin.Context, tx *gorm.DB, req *po.HttpRule) (*po.HttpRule, error) {
	res := &po.HttpRule{}
	err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *ProtocolRuleOperator) SaveHttpRule(ctx *gin.Context, tx *gorm.DB, req *po.HttpRule) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Save(req).Error; err != nil {
		return err
	}
	return nil
}

func (o *ProtocolRuleOperator) ListHttpRulesByServiceId(ctx *gin.Context, tx *gorm.DB, serviceId int64) ([]po.HttpRule, int64, error) {
	var httpRules []po.HttpRule
	var count int64
	var err error
	query := tx.SetCtx(utils.GetGinTraceContext(ctx)).Table((&po.HttpRule{}).TableName()).Where("service_id=?", serviceId)

	err = query.Order("id desc").Find(&httpRules).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, nil
	}

	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return httpRules, count, nil
}

func (o *ProtocolRuleOperator) FindTcpRule(ctx *gin.Context, tx *gorm.DB, req *po.TcpRule) (*po.TcpRule, error) {
	res := &po.TcpRule{}
	err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *ProtocolRuleOperator) SaveTcpRule(ctx *gin.Context, tx *gorm.DB, req *po.TcpRule) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Save(req).Error; err != nil {
		return err
	}
	return nil
}

func (o *ProtocolRuleOperator) ListTcpRulesByServiceId(ctx *gin.Context, tx *gorm.DB, serviceId int64) ([]po.TcpRule, int64, error) {
	var tcpRules []po.TcpRule
	var count int64
	var err error
	query := tx.SetCtx(utils.GetGinTraceContext(ctx)).Table((&po.TcpRule{}).TableName()).Where("service_id=?", serviceId)

	err = query.Order("id desc").Find(&tcpRules).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return tcpRules, count, nil
}

func (o *ProtocolRuleOperator) FindGrpcRule(ctx *gin.Context, tx *gorm.DB, req *po.GrpcRule) (*po.GrpcRule, error) {
	res := &po.GrpcRule{}
	err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Where(req).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *ProtocolRuleOperator) SaveGrpcRule(ctx *gin.Context, tx *gorm.DB, req *po.GrpcRule) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(ctx)).Save(req).Error; err != nil {
		return err
	}
	return nil
}

func (o *ProtocolRuleOperator) ListGrpcRulesByServiceId(ctx *gin.Context, tx *gorm.DB, serviceId int64) ([]po.GrpcRule, int64, error) {
	var grpcRules []po.GrpcRule
	var count int64
	var err error
	query := tx.SetCtx(utils.GetGinTraceContext(ctx)).Table((&po.GrpcRule{}).TableName()).Where("service_id=?", serviceId)

	err = query.Order("id desc").Find(&grpcRules).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, nil
	}

	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return grpcRules, count, nil
}
