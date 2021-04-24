package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/models/dto"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

var jwtService *JwtService

type JwtService struct {
}

func NewJwtService() *JwtService {
	service := &JwtService{}
	return service
}

func GetJwtService() *JwtService {
	if jwtService == nil {
		jwtService = NewJwtService()
	}
	return jwtService
}

func (s *JwtService) GenerateJwt(ctx *gin.Context, tx *gorm.DB, req *dto.JwtReq, appId, secret string) (*dto.JwtRes, error) {
	apps := GetAppService().ListAppsInMemory()

	for _, app := range apps {
		if app.AppId == appId && app.Secret == secret {
			claims := jwt.StandardClaims{
				Issuer:    app.AppId,
				ExpiresAt: time.Now().Add(constants.JwtExpires * time.Second).In(lib.TimeLocation).Unix(),
			}
			token, err := utils.EncodeJwt(claims)
			if err != nil {
				return nil, err
			}

			if req.Type == "" {
				req.Type = constants.JwtType
			}
			if req.Permission == "" {
				req.Permission = constants.JwtReadWrite
			}

			return &dto.JwtRes{
				Token:      token,
				Type:       req.Type,
				ExpireAt:   constants.JwtExpires,
				Permission: req.Permission,
			}, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("secret %s for app %s is incorrect", secret, appId))
}

func (s *JwtService) HttpVerifyJwt(ctx *gin.Context, svc *po.ServiceDetail, tokenString string) error {
	isMatched := false

	if tokenString != "" {
		// verify expire at
		claims, err := utils.DecodeJwt(tokenString)
		if err != nil {
			return err
		}

		// verify issuer
		apps := GetAppService().ListAppsInMemory()
		for _, app := range apps {
			if app.AppId == claims.Issuer {
				ctx.Set("app", app)
				isMatched = true
				break
			}
		}
	}

	if svc.AccessControl.OpenAuth == constants.Enable && !isMatched {
		return errors.New("failed to verify jwt, err: no matched valid app")
	}

	return nil
}

func (s *JwtService) GrpcVerifyJwt(ctx metadata.MD, svc *po.ServiceDetail, tokenString string) error {
	isMatched := false

	if tokenString != "" {
		// verify expire at
		claims, err := utils.DecodeJwt(tokenString)
		if err != nil {
			return err
		}

		// verify issuer
		apps := GetAppService().ListAppsInMemory()
		for _, app := range apps {
			if app.AppId == claims.Issuer {
				ctx.Set("app", utils.Obj2Json(app))
				isMatched = true
				break
			}
		}
	}

	if svc.AccessControl.OpenAuth == constants.Enable && !isMatched {
		return errors.New("failed to verify jwt, err: no matched valid app")
	}

	return nil
}
