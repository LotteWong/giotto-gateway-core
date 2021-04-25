package service

import (
	"fmt"
	"log"
	"net/http/httptest"
	"sync"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/dao/mysql"
	"github.com/LotteWong/giotto-gateway-core/dao/redis"
	"github.com/LotteWong/giotto-gateway-core/models/dto"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/utils"
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var appService *AppService

type AppService struct {
	AppMap   map[string]*po.App
	AppSlice []*po.App
	RWLock   sync.RWMutex
	DCLock   sync.Once
	InitErr  error

	appOperator  *mysql.AppOperator
	appRedisConn *redis.AppOperator
}

func NewAppService() *AppService {
	service := &AppService{
		AppMap:       map[string]*po.App{},
		AppSlice:     []*po.App{},
		RWLock:       sync.RWMutex{},
		DCLock:       sync.Once{},
		InitErr:      nil,
		appOperator:  mysql.NewAppOperator(),
		appRedisConn: redis.NewAppOperator(),
	}
	return service
}

func GetAppService() *AppService {
	if appService == nil {
		appService = NewAppService()
	}
	return appService
}

func (s *AppService) LoadAppsFromRedis() error {
	s.DCLock.Do(func() {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx, err := lib.GetGormPool("default")
		if err != nil {
			s.InitErr = err
			return
		}

		_, apps, err := s.appOperator.FuzzySearchAndPage(ctx, tx, "", 0, 0)
		if err != nil {
			s.InitErr = err
			return
		}

		s.RWLock.Lock()
		defer s.RWLock.Unlock()

		for _, app := range apps {
			tmp := app
			err := s.appRedisConn.SetApp(tmp.Id, &tmp)
			if err != nil {
				s.InitErr = err
				return
			}
		}
	})

	return s.InitErr
}

func (s *AppService) ListAppsFromRedis() ([]*po.App, error) {
	return s.appRedisConn.ListApps()
}

func (s *AppService) LoadAppsIntoMemory() error {
	s.DCLock.Do(func() {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx, err := lib.GetGormPool("default")
		if err != nil {
			s.InitErr = err
			return
		}

		_, apps, err := s.appOperator.FuzzySearchAndPage(ctx, tx, "", 0, 0)
		if err != nil {
			s.InitErr = err
			return
		}

		s.RWLock.Lock()
		defer s.RWLock.Unlock()

		for _, app := range apps {
			tmp := app
			s.AppMap[tmp.AppId] = &tmp
			s.AppSlice = append(s.AppSlice, &tmp)
		}
	})

	return s.InitErr
}

func (s *AppService) ListAppsInMemory() []*po.App {
	return s.AppSlice
}

func (s *AppService) ListApps(ctx *gin.Context, tx *gorm.DB, req *dto.ListAppsReq) (int64, []dto.ListAppItem, error) {
	total, items, err := s.appOperator.FuzzySearchAndPage(ctx, tx, req.Keyword, req.PageIndex, req.PageSize)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("failed to page apps with condition %v, err: %v", req, err))
	}

	var appItems []dto.ListAppItem
	for _, item := range items {
		count, err := GetFlowCountService().GetFlowCount(constants.AppFlowCountPrefix + item.AppId)
		if err != nil {
			return 0, nil, errors.New(fmt.Sprintf("failed to get app flow count of %s, err: %v", item.AppId, err))
		}

		appItem := dto.ListAppItem{
			Id:       item.Id,
			AppId:    item.AppId,
			AppName:  item.AppName,
			Secret:   item.Secret,
			WhiteIps: item.WhiteIps,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
			RealQpd:  count.TotalCount,
			RealQps:  count.Qps,
		}
		appItems = append(appItems, appItem)
	}

	return total, appItems, nil
}

func (s *AppService) ShowApp(ctx *gin.Context, tx *gorm.DB, appId int64) (*po.App, error) {
	var app *po.App
	var redisErr, mysqlErr error

	app, redisErr = s.appRedisConn.GetApp(appId)
	if redisErr != nil {
		app, mysqlErr = s.appOperator.Find(ctx, tx, &po.App{Id: appId})
		if mysqlErr != nil {
			return nil, mysqlErr
		}

		log.Printf("show app %d miss redis, query in mysql\n", appId)

		redisErr = s.appRedisConn.SetApp(appId, app)
		if redisErr != nil {
			return nil, redisErr
		}
	}

	return app, nil
}

func (s *AppService) CreateApp(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateAppReq) (*po.App, error) {
	err := s.validCreateApp(ctx, tx, req)
	if err != nil {
		return nil, err
	}

	tx = tx.Begin()

	app := &po.App{
		AppId:    req.AppId,
		AppName:  req.AppName,
		Secret:   req.Secret,
		WhiteIps: req.WhiteIps,
		Qps:      req.Qps,
		Qpd:      req.Qpd,
	}
	if err := s.appOperator.Save(ctx, tx, app); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	if err := s.appRedisConn.SetApp(app.Id, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *AppService) UpdateApp(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateAppReq, appId int64) (*po.App, error) {
	err := s.validUpdateApp(req)
	if err != nil {
		return nil, err
	}

	tx = tx.Begin()

	app, err := s.appOperator.Find(ctx, tx, &po.App{Id: appId})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	app.AppName = req.AppName
	app.Secret = req.Secret
	app.WhiteIps = req.WhiteIps
	app.Qpd = req.Qpd
	app.Qps = req.Qps
	if err := s.appOperator.Save(ctx, tx, app); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	if err := s.appRedisConn.SetApp(appId, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *AppService) DeleteApp(ctx *gin.Context, tx *gorm.DB, appId int64) error {
	app, err := s.appOperator.Find(ctx, tx, &po.App{Id: appId})
	if err != nil {
		return err
	}

	app.IsDelete = 1
	err = s.appOperator.Save(ctx, tx, app)
	if err != nil {
		return err
	}

	if err := s.appRedisConn.DelApp(app.Id); err != nil {
		return err
	}

	return nil
}

func (s *AppService) validCreateApp(ctx *gin.Context, tx *gorm.DB, req *dto.CreateOrUpdateAppReq) error {
	// check whether app id is duplicated
	if _, err := s.appOperator.Find(ctx, tx, &po.App{AppId: req.AppId}); err == nil {
		return errors.New(fmt.Sprintf("app id %s is duplicated", req.AppId))
	}

	// generate random secret if not specified
	if req.Secret == "" {
		secret, err := utils.MD5(req.AppId)
		if err != nil {
			return err
		}
		req.Secret = secret
	}

	return nil
}

func (s *AppService) validUpdateApp(req *dto.CreateOrUpdateAppReq) error {
	// generate random secret if not specified
	if req.Secret == "" {
		secret, err := utils.MD5(req.AppId)
		if err != nil {
			return err
		}
		req.Secret = secret
	}

	return nil
}
