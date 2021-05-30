package service

import (
	"strings"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var ipService *IpService

type IpService struct {
}

func NewIpService() *IpService {
	service := &IpService{}
	return service
}

func GetIpService() *IpService {
	if ipService == nil {
		ipService = NewIpService()
	}
	return ipService
}

func (s *IpService) VerifyIpListForService(ctx *gin.Context, svc *po.ServiceDetail) error {
	var whiteIpList []string
	var blackIpList []string
	openAuth := svc.AccessControl.OpenAuth
	whiteList := svc.AccessControl.WhiteList
	if whiteList != "" {
		whiteIpList = strings.Split(whiteList, ",")
	}
	blackList := svc.AccessControl.BlackList
	if blackList != "" {
		blackIpList = strings.Split(blackList, ",")
	}

	if openAuth == constants.Enable {
		if len(whiteIpList) > 0 { // white list has higher priority
			if !checkStrInSlice(whiteIpList, ctx.ClientIP()) {
				return errors.Errorf("ip %s is not in service %s's white ip list", ctx.ClientIP(), svc.Info.ServiceName)
			}
		} else { // black list has lower priority
			if len(blackIpList) > 0 {
				if checkStrInSlice(blackIpList, ctx.ClientIP()) {
					return errors.Errorf("ip %s is in service %s's black ip list", ctx.ClientIP(), svc.Info.ServiceName)
				}
			}
		}
	}
	return nil
}

func (s *IpService) VerifyIpListForApp(ctx *gin.Context, app *po.App) error {
	var whiteIpList []string
	var blackIpList []string
	whiteIps := app.WhiteIps
	if whiteIps != "" {
		whiteIpList = strings.Split(whiteIps, ",")
	}
	blackIps := app.BlackIps
	if blackIps != "" {
		blackIpList = strings.Split(blackIps, ",")
	}

	if len(whiteIpList) > 0 { // white list has higher priority
		if !checkStrInSlice(whiteIpList, ctx.ClientIP()) {
			return errors.Errorf("ip %s is not in app %s's white ip list", ctx.ClientIP(), app.AppId)
		}
	} else { // black list has lower priority
		if len(blackIpList) > 0 {
			if checkStrInSlice(blackIpList, ctx.ClientIP()) {
				return errors.Errorf("ip %s is in app %s's black ip list", ctx.ClientIP(), app.AppId)
			}
		}
	}
	return nil
}

func checkStrInSlice(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
