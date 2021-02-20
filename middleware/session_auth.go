package middleware

import (
	"errors"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if loginSession, ok := session.Get(constants.LoginSessionKey).(string); !ok || loginSession == "" {
			ResponseError(c, InternalErrorCode, errors.New("user not login"))
			c.Abort()
			return
		}
		c.Next()
	}
}
