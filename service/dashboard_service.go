package service

import (
	"fmt"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/dao"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var dashboardService *DashboardService

type DashboardService struct {
	serviceOperator *dao.ServiceOperator
}

func NewDashboardService() *DashboardService {
	service := &DashboardService{
		serviceOperator: dao.NewServiceOperator(),
	}
	return service
}

func GetDashboardService() *DashboardService {
	if dashboardService == nil {
		dashboardService = NewDashboardService()
	}
	return dashboardService
}

func (s *DashboardService) GetServicePercentage(ctx *gin.Context, tx *gorm.DB) (*dto.ServicePercentageItems, error) {
	groups, err := s.serviceOperator.GroupByServiceType(ctx, tx)
	if err != nil {
		return nil, err
	}

	var legends []string
	var records []dto.ServicePercentageItem
	for _, group := range groups {
		legend, ok := constants.ServiceTypeMap[group.ServiceType]
		if !ok {
			return nil, errors.New(fmt.Sprintf("service type %d not found", group.ServiceType))
		}
		legends = append(legends, legend)

		record := &dto.ServicePercentageItem{
			ServiceLegend: legend,
			ServiceType:   group.ServiceType,
			ServiceCount:  group.ServiceCount,
		}
		records = append(records, *record)
	}

	return &dto.ServicePercentageItems{
		Legends: legends,
		Records: records,
	}, nil
}
