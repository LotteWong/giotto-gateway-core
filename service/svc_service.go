package service

import (
	"fmt"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/dao"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var svcService *SvcService

type SvcService struct {
	serviceOperator       *dao.ServiceOperator
	protocolRuleOperator  *dao.ProtocolRuleOperator
	loadBalanceOperator   *dao.LoadBalanceOperator
	accessControlOperator *dao.AccessControlOperator
}

func NewSvcService() *SvcService {
	service := &SvcService{
		serviceOperator:       dao.NewServiceOperator(),
		protocolRuleOperator:  dao.NewProtocolRuleOperator(),
		loadBalanceOperator:   dao.NewLoadBalanceOperator(),
		accessControlOperator: dao.NewAccessControlOperator(),
	}
	return service
}

func GetSvcService() *SvcService {
	if svcService == nil {
		svcService = NewSvcService()
	}
	return svcService
}

func (s *SvcService) ListServices(ctx *gin.Context, tx *gorm.DB, req *dto.ListServicesReq) (int64, []dto.Service, error) {
	total, serviceInfoItems, err := s.serviceOperator.FuzzySearchAndPage(ctx, tx, req.Keyword, req.PageIndex, req.PageSize)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("failed to page services with condition %v, err: %v", req, err))
	}

	var serviceItems []dto.Service
	for _, serviceInfoItem := range serviceInfoItems {
		serviceDetail, err := s.getServiceDetail(ctx, tx, serviceInfoItem.Id)
		if err != nil {
			return 0, nil, errors.New(fmt.Sprintf("failed to get service detail of %s, err: %v", serviceInfoItem.ServiceName, err))
		}

		serviceItem := dto.Service{
			Id:          serviceInfoItem.Id,
			ServiceName: serviceInfoItem.ServiceName,
			ServiceDesc: serviceInfoItem.ServiceDesc,
			ServiceAddr: s.concatServiceAddr(serviceDetail),
			Qps:         0, // TODO
			Qpd:         0, // TODO
			TotalNode:   len(serviceDetail.LoadBalance.GetEnabledIpList()),
		}
		serviceItems = append(serviceItems, serviceItem)
	}

	return total, serviceItems, nil
}

func (s *SvcService) getServiceDetail(ctx *gin.Context, tx *gorm.DB, serviceId int64) (*po.ServiceDetail, error) {
	var err error

	serviceInfo := &po.ServiceInfo{Id: serviceId}
	serviceInfo, err = s.serviceOperator.Find(ctx, tx, serviceInfo)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	httpRule := &po.HttpRule{ServiceId: serviceId}
	httpRule, err = s.protocolRuleOperator.FindHttpRule(ctx, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule := &po.TcpRule{ServiceId: serviceId}
	tcpRule, err = s.protocolRuleOperator.FindTcpRule(ctx, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule := &po.GrpcRule{ServiceId: serviceId}
	grpcRule, err = s.protocolRuleOperator.FindGrpcRule(ctx, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance := &po.LoadBalance{ServiceId: serviceId}
	loadBalance, err = s.loadBalanceOperator.Find(ctx, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	accessControl := &po.AccessControl{ServiceId: serviceId}
	accessControl, err = s.accessControlOperator.Find(ctx, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &po.ServiceDetail{
		Info:          serviceInfo,
		TcpRule:       tcpRule,
		HttpRule:      httpRule,
		GrpcRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}, nil
}

func (s *SvcService) concatServiceAddr(detail *po.ServiceDetail) string {
	addr := "unknown"
	clusterIp := lib.GetStringConf("base.cluster.cluster_ip")
	clusterPort := lib.GetStringConf("base.cluster.cluster_port")
	clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

	if detail.HttpRule != nil {
		if detail.Info.ServiceType == constants.ServiceTypeHttp && detail.HttpRule.RuleType == constants.HttpRuleTypePrefixUrl && detail.HttpRule.NeedHttps == constants.DisableHttps {
			addr = fmt.Sprintf("%s:%s%s", clusterIp, clusterPort, detail.HttpRule.Rule)
			return addr
		}

		if detail.Info.ServiceType == constants.ServiceTypeHttp && detail.HttpRule.RuleType == constants.HttpRuleTypePrefixUrl && detail.HttpRule.NeedHttps == constants.EnableHttps {
			addr = fmt.Sprintf("%s:%s%s", clusterIp, clusterSSLPort, detail.HttpRule.Rule)
			return addr
		}

		if detail.Info.ServiceType == constants.ServiceTypeHttp && detail.HttpRule.RuleType == constants.HttpRuleTypeDomain {
			addr = detail.HttpRule.Rule
			return addr
		}
	}

	if detail.TcpRule != nil {
		if detail.Info.ServiceType == constants.ServiceTypeTcp {
			addr = fmt.Sprintf("%s:%d", clusterIp, detail.TcpRule.Port)
			return addr
		}
	}

	if detail.GrpcRule != nil {
		if detail.Info.ServiceType == constants.ServiceTypeGrpc {
			addr = fmt.Sprintf("%s:%d", clusterIp, detail.GrpcRule.Port)
			return addr
		}
	}

	return addr
}

func (s *SvcService) DeleteServices(ctx *gin.Context, tx *gorm.DB, serviceId int64) error {
	serviceInfo, err := s.serviceOperator.Find(ctx, tx, &po.ServiceInfo{Id: serviceId})
	if err != nil {
		return err
	}

	serviceInfo.IsDelete = 1
	err = s.serviceOperator.Save(ctx, tx, serviceInfo)
	if err != nil {
		return err
	}

	return nil
}
