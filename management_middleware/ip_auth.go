package management_middleware

import (
	"fmt"
	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func IpAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var isMatched bool
		var authMode string
		allowIps := lib.GetStringSliceConf("base.http.allow_ips")
		blockIps := lib.GetStringSliceConf("base.http.block_ips")

		if len(allowIps) > 0 {
			authMode = "allow"
			isMatched = false
			for _, host := range allowIps {
				if c.ClientIP() == host {
					isMatched = true
					break
				}
			}
		} else {
			authMode = "block"
			isMatched = true
			if len(blockIps) > 0 {
				for _, host := range blockIps {
					if c.ClientIP() == host {
						isMatched = false
						break
					}
				}
			}
		}

		if !isMatched {
			switch authMode {
			case "allow":
				common_middleware.ResponseError(c, common_middleware.InternalErrorCode, errors.New(fmt.Sprintf("ip %s not in allow ip list", c.ClientIP())))
			case "block":
				common_middleware.ResponseError(c, common_middleware.InternalErrorCode, errors.New(fmt.Sprintf("ip %s is in block ip list", c.ClientIP())))
			}
			c.Abort()
			return
		}
		c.Next()
	}
}
