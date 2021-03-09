package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/dao"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var userService *UserService

type UserService struct {
	userOperator *dao.UserOperator
}

func NewUserService() *UserService {
	service := &UserService{
		userOperator: dao.NewUserOperator(),
	}
	return service
}

func GetUserService() *UserService {
	if userService == nil {
		userService = NewUserService()
	}
	return userService
}

func (s *UserService) GetUserInfo(ctx *gin.Context) (*dto.UserInfo, error) {
	// get user info from session
	session := sessions.Default(ctx)
	loginSessionStruct := &dto.LoginSession{}
	loginSessionBytes := []byte(session.Get(constants.LoginSessionKey).(string))
	if err := json.Unmarshal(loginSessionBytes, loginSessionStruct); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to unmarshal data %v, err: %v", loginSessionBytes, err))
	}

	// get user info by default
	userInfo := &dto.UserInfo{
		ID:       loginSessionStruct.Id,
		Username: loginSessionStruct.Username,
		LoginAt:  loginSessionStruct.LoginAt,
		Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Intro:    "I am the administrator.",
		Roles:    []string{"admin"},
	}

	return userInfo, nil
}

func (s *UserService) ChangeUserPassword(ctx *gin.Context, tx *gorm.DB, req *dto.ChangeUserPwdReq) error {
	// check login info
	session := sessions.Default(ctx)
	loginSessionStruct := &dto.LoginSession{}
	loginSessionBytes := []byte(session.Get(constants.LoginSessionKey).(string))
	if err := json.Unmarshal(loginSessionBytes, loginSessionStruct); err != nil {
		return errors.New(fmt.Sprintf("failed to unmarshal data %v, err: %v", loginSessionBytes, err))
	}

	// check user info
	user := &po.Admin{
		Id:       loginSessionStruct.Id,
		Username: loginSessionStruct.Username,
		IsDelete: 0,
	}
	user, err := s.userOperator.Find(ctx, tx, user)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to find user %s, err: %v", loginSessionStruct.Username, err))
	}

	// commit password change
	saltPassword := utils.GenSaltPwd(user.Salt, req.Password)
	user.Password = saltPassword
	if err := s.userOperator.Save(ctx, tx, user); err != nil {
		return errors.New(fmt.Sprintf("failed to save user %s, err: %v", loginSessionStruct.Username, err))
	}

	return nil
}
