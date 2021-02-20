package service

import (
	"encoding/json"
	"fmt"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/dao"
	"github.com/LotteWong/giotto-gateway/dto"
	"github.com/LotteWong/giotto-gateway/po"
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

var loginService *LoginService

type LoginService struct {
	userOperator *dao.UserOperator
}

func NewLoginService() *LoginService {
	service := &LoginService{
		userOperator: dao.NewUserOperator(),
	}
	return service
}

func GetLoginService() *LoginService {
	if loginService == nil {
		loginService = NewLoginService()
	}
	return loginService
}

func (s *LoginService) Login(ctx *gin.Context, tx *gorm.DB, req *dto.LoginReq) (*po.Admin, error) {
	// get user info
	user := &po.Admin{}
	user, err := s.userOperator.Find(ctx, tx, &po.Admin{
		Username: req.Username,
		IsDelete: 0,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to find user %s, err: %v", req.Username, err))
	}

	// encrypt and compare
	saltPassword := utils.GenSaltPwd(user.Salt, req.Password)
	if user.Password != saltPassword {
		return nil, errors.New(fmt.Sprintf("password %s for user %s is incorrect", req.Password, req.Username))
	}

	// save login session
	loginSessionStruct := &dto.LoginSession{
		Id:       user.Id,
		Username: user.Username,
		LoginAt:  time.Now(),
	}
	loginSessionBytes, err := json.Marshal(loginSessionStruct)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to marshal data %v, err: %v", loginSessionStruct, err))
	}
	loginSession := sessions.Default(ctx)
	loginSession.Set(constants.LoginSessionKey, string(loginSessionBytes))
	if err = loginSession.Save(); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to set session value %v, err: %v", loginSessionBytes, err))
	}

	return user, nil
}

func (s *LoginService) Logout(ctx *gin.Context) error {
	session := sessions.Default(ctx)
	session.Delete(constants.LoginSessionKey)
	if err := session.Save(); err != nil {
		return errors.New(fmt.Sprintf("fail to delete session key %s, err: %v", constants.LoginSessionKey, err))
	}
	return nil
}
