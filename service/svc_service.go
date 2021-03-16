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
	"strings"
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

func (s *SvcService) ListServices(ctx *gin.Context, tx *gorm.DB, req *dto.ListServicesReq) (int64, []dto.ListServiceItem, error) {
	total, serviceInfoItems, err := s.serviceOperator.FuzzySearchAndPage(ctx, tx, req.Keyword, req.PageIndex, req.PageSize)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("failed to page services with condition %v, err: %v", req, err))
	}

	var serviceItems []dto.ListServiceItem
	for _, serviceInfoItem := range serviceInfoItems {
		serviceDetail, err := s.getServiceDetail(ctx, tx, serviceInfoItem.Id)
		if err != nil {
			return 0, nil, errors.New(fmt.Sprintf("failed to get service detail of %s, err: %v", serviceInfoItem.ServiceName, err))
		}

		serviceItem := dto.ListServiceItem{
			Id:          serviceInfoItem.Id,
			ServiceName: serviceInfoItem.ServiceName,
			ServiceDesc: serviceInfoItem.ServiceDesc,
			ServiceType: serviceInfoItem.ServiceType,
			ServiceAddr: s.concatServiceAddr(serviceDetail),
			Qps:         0, // TODO
			Qpd:         0, // TODO
			TotalNode:   len(serviceDetail.LoadBalance.GetEnabledIpList()),
		}
		serviceItems = append(serviceItems, serviceItem)
	}

	return total, serviceItems, nil
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

func (s *SvcService) ShowService(ctx *gin.Context, tx *gorm.DB, serviceId int64) (*po.ServiceDetail, error) {
	serviceDetail, err := s.getServiceDetail(ctx, tx, serviceId)
	if err != nil {
		return nil, err
	}
	return serviceDetail, nil
}

func (s *SvcService) CreateHttpService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateHttpServiceReq) (*po.ServiceDetail, error) {
	err := s.validCreateHttpService(ctx, tx, req)
	if err != nil {
		return nil, err
	}

	tx = tx.Begin()

	// save into service info table
	serviceInfo := &po.ServiceInfo{
		ServiceName: req.ServiceName,
		ServiceDesc: req.ServiceDesc,
		ServiceType: constants.ServiceTypeHttp,
	}
	if err := s.serviceOperator.Save(ctx, tx, serviceInfo); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into http rule table
	httpRule := &po.HttpRule{
		ServiceId:       serviceInfo.Id,
		RuleType:        req.RuleType,
		Rule:            req.Rule,
		NeedHttps:       req.NeedHttps,
		NeedStripUri:    req.NeedStripUri,
		NeedWebsocket:   req.NeedWebsocket,
		UrlRewrite:      req.UrlRewrite,
		HeaderTransform: req.HeaderTransform,
	}
	if err := s.protocolRuleOperator.SaveHttpRule(ctx, tx, httpRule); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into load balance table
	loadBalance := &po.LoadBalance{
		ServiceId:              serviceInfo.Id,
		RoundType:              req.RoundType,
		IpList:                 req.IpList,
		WeightList:             req.WeightList,
		UpstreamConnectTimeout: req.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  req.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    req.UpstreamIdleTimeout,
		UpstreamMaxIdle:        req.UpstreamMaxIdle,
	}
	if err := s.loadBalanceOperator.Save(ctx, tx, loadBalance); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into access control table
	accessControl := &po.AccessControl{
		ServiceId:            serviceInfo.Id,
		OpenAuth:             req.OpenAuth,
		BlackList:            req.BlackList,
		WhiteList:            req.WhiteList,
		ClientIpFlowLimit:    req.ClientIpFlowLimit,
		ServiceHostFlowLimit: req.ServiceHostFlowLimit,
	}
	if err := s.accessControlOperator.Save(ctx, tx, accessControl); err != nil {
		tx.Rollback()
		return nil, err
	}

	serviceDetail, err := s.getServiceDetail(ctx, tx, serviceInfo.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return serviceDetail, nil
}

func (s *SvcService) UpdateHttpService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateHttpServiceReq, serviceId int64) (*po.ServiceDetail, error) {
	err := s.validUpdateHttpService(ctx, tx, req, serviceId)
	if err != nil {
		return nil, err
	}

	tx = tx.Begin()

	serviceDetail, err := s.getServiceDetail(ctx, tx, serviceId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into service info table
	serviceInfo := serviceDetail.Info
	serviceInfo.ServiceDesc = req.ServiceDesc
	if err := s.serviceOperator.Save(ctx, tx, serviceInfo); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into http rule table
	httpRule := serviceDetail.HttpRule
	httpRule.NeedHttps = req.NeedHttps
	httpRule.NeedStripUri = req.NeedStripUri
	httpRule.NeedWebsocket = req.NeedWebsocket
	httpRule.UrlRewrite = req.UrlRewrite
	httpRule.HeaderTransform = req.HeaderTransform
	if err := s.protocolRuleOperator.SaveHttpRule(ctx, tx, httpRule); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into load balance table
	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = req.RoundType
	loadBalance.IpList = req.IpList
	loadBalance.WeightList = req.WeightList
	loadBalance.UpstreamConnectTimeout = req.UpstreamConnectTimeout
	loadBalance.UpstreamHeaderTimeout = req.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = req.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = req.UpstreamMaxIdle
	if err := s.loadBalanceOperator.Save(ctx, tx, loadBalance); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into access control table
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = req.OpenAuth
	accessControl.BlackList = req.BlackList
	accessControl.WhiteList = req.WhiteList
	accessControl.ClientIpFlowLimit = req.ClientIpFlowLimit
	accessControl.ServiceHostFlowLimit = req.ServiceHostFlowLimit
	if err := s.accessControlOperator.Save(ctx, tx, accessControl); err != nil {
		tx.Rollback()
		return nil, err
	}

	serviceDetail, err = s.getServiceDetail(ctx, tx, serviceInfo.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return serviceDetail, nil
}

func (s *SvcService) CreateTcpService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateTcpServiceReq) (*po.ServiceDetail, error) {
	err := s.validCreateTcpService(ctx, tx, req)
	if err != nil {
		return nil, err
	}

	tx = tx.Begin()

	// save into service info table
	serviceInfo := &po.ServiceInfo{
		ServiceName: req.ServiceName,
		ServiceDesc: req.ServiceDesc,
		ServiceType: constants.ServiceTypeTcp,
	}
	if err := s.serviceOperator.Save(ctx, tx, serviceInfo); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into tcp rule table
	tcpRule := &po.TcpRule{
		ServiceId: serviceInfo.Id,
		Port:      req.Port,
	}
	if err := s.protocolRuleOperator.SaveTcpRule(ctx, tx, tcpRule); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into load balance table
	loadBalance := &po.LoadBalance{
		ServiceId:  serviceInfo.Id,
		RoundType:  req.RoundType,
		IpList:     req.IpList,
		WeightList: req.WeightList,
		ForbidList: req.ForbidList,
	}
	if err := s.loadBalanceOperator.Save(ctx, tx, loadBalance); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into access control table
	accessControl := &po.AccessControl{
		ServiceId:            serviceInfo.Id,
		OpenAuth:             req.OpenAuth,
		BlackList:            req.BlackList,
		WhiteList:            req.WhiteList,
		WhiteHostName:        req.WhiteHostName,
		ClientIpFlowLimit:    req.ClientIpFlowLimit,
		ServiceHostFlowLimit: req.ServiceHostFlowLimit,
	}
	if err := s.accessControlOperator.Save(ctx, tx, accessControl); err != nil {
		tx.Rollback()
		return nil, err
	}

	serviceDetail, err := s.getServiceDetail(ctx, tx, serviceInfo.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return serviceDetail, nil
}

func (s *SvcService) UpdateTcpService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateTcpServiceReq, serviceId int64) (*po.ServiceDetail, error) {
	err := s.validUpdateTcpService(ctx, tx, req, serviceId)
	if err != nil {
		return nil, err
	}

	tx = tx.Begin()

	serviceDetail, err := s.getServiceDetail(ctx, tx, serviceId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into service info table
	serviceInfo := serviceDetail.Info
	serviceInfo.ServiceDesc = req.ServiceDesc
	if err := s.serviceOperator.Save(ctx, tx, serviceInfo); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into load balance table
	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = req.RoundType
	loadBalance.IpList = req.IpList
	loadBalance.WeightList = req.WeightList
	loadBalance.ForbidList = req.ForbidList
	if err := s.loadBalanceOperator.Save(ctx, tx, loadBalance); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into access control table
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = req.OpenAuth
	accessControl.BlackList = req.BlackList
	accessControl.WhiteList = req.WhiteList
	accessControl.WhiteHostName = req.WhiteHostName
	accessControl.ClientIpFlowLimit = req.ClientIpFlowLimit
	accessControl.ServiceHostFlowLimit = req.ServiceHostFlowLimit
	if err := s.accessControlOperator.Save(ctx, tx, accessControl); err != nil {
		tx.Rollback()
		return nil, err
	}

	serviceDetail, err = s.getServiceDetail(ctx, tx, serviceInfo.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return serviceDetail, nil
}

func (s *SvcService) CreateGrpcService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateGrpcServiceReq) (*po.ServiceDetail, error) {
	err := s.validCreateGrpcService(ctx, tx, req)
	if err != nil {
		return nil, err
	}

	tx = tx.Begin()

	// save into service info table
	serviceInfo := &po.ServiceInfo{
		ServiceName: req.ServiceName,
		ServiceDesc: req.ServiceDesc,
		ServiceType: constants.ServiceTypeGrpc,
	}
	if err := s.serviceOperator.Save(ctx, tx, serviceInfo); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into tcp rule table
	grpcRule := &po.GrpcRule{
		ServiceId:       serviceInfo.Id,
		Port:            req.Port,
		HeaderTransform: req.HeaderTransform,
	}
	if err := s.protocolRuleOperator.SaveGrpcRule(ctx, tx, grpcRule); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into load balance table
	loadBalance := &po.LoadBalance{
		ServiceId:  serviceInfo.Id,
		RoundType:  req.RoundType,
		IpList:     req.IpList,
		WeightList: req.WeightList,
		ForbidList: req.ForbidList,
	}
	if err := s.loadBalanceOperator.Save(ctx, tx, loadBalance); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into access control table
	accessControl := &po.AccessControl{
		ServiceId:            serviceInfo.Id,
		OpenAuth:             req.OpenAuth,
		BlackList:            req.BlackList,
		WhiteList:            req.WhiteList,
		WhiteHostName:        req.WhiteHostName,
		ClientIpFlowLimit:    req.ClientIpFlowLimit,
		ServiceHostFlowLimit: req.ServiceHostFlowLimit,
	}
	if err := s.accessControlOperator.Save(ctx, tx, accessControl); err != nil {
		tx.Rollback()
		return nil, err
	}

	serviceDetail, err := s.getServiceDetail(ctx, tx, serviceInfo.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return serviceDetail, nil
}

func (s *SvcService) UpdateGrpcService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateGrpcServiceReq, serviceId int64) (*po.ServiceDetail, error) {
	err := s.validUpdateGrpcService(ctx, tx, req, serviceId)
	if err != nil {
		return nil, err
	}

	tx = tx.Begin()

	serviceDetail, err := s.getServiceDetail(ctx, tx, serviceId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into service info table
	serviceInfo := serviceDetail.Info
	serviceInfo.ServiceDesc = req.ServiceDesc
	if err := s.serviceOperator.Save(ctx, tx, serviceInfo); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into http rule table
	grpcRule := serviceDetail.GrpcRule
	grpcRule.HeaderTransform = req.HeaderTransform
	if err := s.protocolRuleOperator.SaveGrpcRule(ctx, tx, grpcRule); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into load balance table
	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = req.RoundType
	loadBalance.IpList = req.IpList
	loadBalance.WeightList = req.WeightList
	loadBalance.ForbidList = req.ForbidList
	if err := s.loadBalanceOperator.Save(ctx, tx, loadBalance); err != nil {
		tx.Rollback()
		return nil, err
	}

	// save into access control table
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = req.OpenAuth
	accessControl.BlackList = req.BlackList
	accessControl.WhiteList = req.WhiteList
	accessControl.WhiteHostName = req.WhiteHostName
	accessControl.ClientIpFlowLimit = req.ClientIpFlowLimit
	accessControl.ServiceHostFlowLimit = req.ServiceHostFlowLimit
	if err := s.accessControlOperator.Save(ctx, tx, accessControl); err != nil {
		tx.Rollback()
		return nil, err
	}

	serviceDetail, err = s.getServiceDetail(ctx, tx, serviceInfo.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return serviceDetail, nil
}

func (s *SvcService) DeleteService(ctx *gin.Context, tx *gorm.DB, serviceId int64) error {
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

func (s *SvcService) getServiceDetail(ctx *gin.Context, tx *gorm.DB, serviceId int64) (*po.ServiceDetail, error) {
	var err error

	serviceInfo := &po.ServiceInfo{Id: serviceId}
	serviceInfo, err = s.serviceOperator.Find(ctx, tx, serviceInfo)
	if err != nil {
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

func (s *SvcService) validCreateHttpService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateHttpServiceReq) error {
	// check whether service name is duplicated
	if _, err := s.serviceOperator.Find(ctx, tx, &po.ServiceInfo{ServiceName: req.ServiceName}); err == nil {
		return errors.New(fmt.Sprintf("service name %s is duplicated", req.ServiceName))
	}

	// check whether http rule is duplicated
	if _, err := s.protocolRuleOperator.FindHttpRule(ctx, tx, &po.HttpRule{RuleType: req.RuleType, Rule: req.Rule}); err == nil {
		return errors.New(fmt.Sprintf("http rule %s is duplicated", req.Rule))
	}

	// check whether ip list is corresponded to weight list
	ipListLen := len(strings.Split(req.IpList, ","))
	weightListLen := len(strings.Split(req.WeightList, ","))
	if ipListLen != weightListLen {
		return errors.New("ip list's length is not corresponded to weight list's length")
	}

	return nil
}

func (s *SvcService) validUpdateHttpService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateHttpServiceReq, serviceId int64) error {
	// check whether service exists
	serviceInfo, err := s.serviceOperator.Find(ctx, tx, &po.ServiceInfo{Id: serviceId})
	if err != nil {
		return errors.New(fmt.Sprintf("service %d not exist, err: %v", serviceId, err))
	}
	if serviceInfo.ServiceType != constants.ServiceTypeHttp {
		return errors.New(fmt.Sprintf("update http service error occurs, can not update service %d of other type", serviceId))
	}

	// check whether ip list is corresponded to weight list
	ipListLen := len(strings.Split(req.IpList, ","))
	weightListLen := len(strings.Split(req.WeightList, ","))
	if ipListLen != weightListLen {
		return errors.New("ip list's length is not corresponded to weight list's length")
	}

	return nil
}

func (s *SvcService) validCreateTcpService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateTcpServiceReq) error {
	// check whether service name is duplicated
	if _, err := s.serviceOperator.Find(ctx, tx, &po.ServiceInfo{ServiceName: req.ServiceName}); err == nil {
		return errors.New(fmt.Sprintf("service name %s is duplicated", req.ServiceName))
	}

	// check whether tcp port is duplicated
	if _, err := s.protocolRuleOperator.FindTcpRule(ctx, tx, &po.TcpRule{Port: req.Port}); err == nil {
		return errors.New(fmt.Sprintf("tcp port %d is duplicated", req.Port))
	}
	if _, err := s.protocolRuleOperator.FindGrpcRule(ctx, tx, &po.GrpcRule{Port: req.Port}); err == nil {
		return errors.New(fmt.Sprintf("grpc port %d is duplicated", req.Port))
	}

	// check whether ip list is corresponded to weight list
	ipListLen := len(strings.Split(req.IpList, ","))
	weightListLen := len(strings.Split(req.WeightList, ","))
	if ipListLen != weightListLen {
		return errors.New("ip list's length is not corresponded to weight list's length")
	}

	return nil
}

func (s *SvcService) validUpdateTcpService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateTcpServiceReq, serviceId int64) error {
	// check whether service exists
	serviceInfo, err := s.serviceOperator.Find(ctx, tx, &po.ServiceInfo{Id: serviceId})
	if err != nil {
		return errors.New(fmt.Sprintf("service %d not exist, err: %v", serviceId, err))
	}
	if serviceInfo.ServiceType != constants.ServiceTypeTcp {
		return errors.New(fmt.Sprintf("update tcp service error occurs, can not update service %d of other type", serviceId))
	}

	// check whether ip list is corresponded to weight list
	ipListLen := len(strings.Split(req.IpList, ","))
	weightListLen := len(strings.Split(req.WeightList, ","))
	if ipListLen != weightListLen {
		return errors.New("ip list's length is not corresponded to weight list's length")
	}

	return nil
}

func (s *SvcService) validCreateGrpcService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateGrpcServiceReq) error {
	// check whether service name is duplicated
	if _, err := s.serviceOperator.Find(ctx, tx, &po.ServiceInfo{ServiceName: req.ServiceName}); err == nil {
		return errors.New(fmt.Sprintf("service name %s is duplicated", req.ServiceName))
	}

	// check whether tcp port is duplicated
	if _, err := s.protocolRuleOperator.FindTcpRule(ctx, tx, &po.TcpRule{Port: req.Port}); err == nil {
		return errors.New(fmt.Sprintf("tcp port %d is duplicated", req.Port))
	}
	if _, err := s.protocolRuleOperator.FindGrpcRule(ctx, tx, &po.GrpcRule{Port: req.Port}); err == nil {
		return errors.New(fmt.Sprintf("grpc port %d is duplicated", req.Port))
	}

	// check whether ip list is corresponded to weight list
	ipListLen := len(strings.Split(req.IpList, ","))
	weightListLen := len(strings.Split(req.WeightList, ","))
	if ipListLen != weightListLen {
		return errors.New("ip list's length is not corresponded to weight list's length")
	}

	return nil
}

func (s *SvcService) validUpdateGrpcService(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateGrpcServiceReq, serviceId int64) error {
	// check whether service exists
	// check whether service exists
	serviceInfo, err := s.serviceOperator.Find(ctx, tx, &po.ServiceInfo{Id: serviceId})
	if err != nil {
		return errors.New(fmt.Sprintf("service %d not exist, err: %v", serviceId, err))
	}
	if serviceInfo.ServiceType != constants.ServiceTypeGrpc {
		return errors.New(fmt.Sprintf("update grpc service error occurs, can not update service %d of other type", serviceId))
	}

	// check whether ip list is corresponded to weight list
	ipListLen := len(strings.Split(req.IpList, ","))
	weightListLen := len(strings.Split(req.WeightList, ","))
	if ipListLen != weightListLen {
		return errors.New("ip list's length is not corresponded to weight list's length")
	}

	return nil
}
