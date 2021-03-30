package http_proxy_middleware

import (
	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func HttpRouteRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpServiceInterface, ok := c.Get("service")
		if !ok {
			common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New("service not found"))
			c.Abort()
			return
		}
		httpServiceDetail := httpServiceInterface.(*po.ServiceDetail)

		// need to rewrite url
		for _, rule := range strings.Split(httpServiceDetail.HttpRule.UrlRewrite, ",") {
			items := strings.Split(rule, " ")
			if len(items) != 2 {
				log.Println("url rewrite format error")
				continue
			}
			ruleBeReplaced, ruleToReplace := items[0], items[1]
			regExp, err := regexp.Compile(ruleBeReplaced)
			if err != nil {
				log.Println(errors.New(err.Error()))
			}
			c.Request.URL.Path = string(regExp.ReplaceAll([]byte(c.Request.URL.Path), []byte(ruleToReplace)))
		}

		// need to strip uri
		if httpServiceDetail.HttpRule.NeedStripUri == constants.Enable && httpServiceDetail.HttpRule.RuleType == constants.HttpRuleTypePrefixUrl {
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, httpServiceDetail.HttpRule.Rule, "", 1)
		}

		// need to transform header
		for _, rule := range strings.Split(httpServiceDetail.HttpRule.HeaderTransform, ",") {
			items := strings.Split(rule, " ")
			if len(items) < 2 || len(items) > 3 {
				log.Println("url rewrite format error")
				continue
			}
			action := items[0]

			if len(items) == 2 {
				key := items[1]
				if action == "delete" {
					c.Request.Header.Del(key)
				}
			} else {
				key, value := items[1], items[2]
				if action == "add" || action == "edit" {
					c.Request.Header.Set(key, value)
				}
			}

		}

		c.Next()
	}
}
