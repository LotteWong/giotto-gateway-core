package common_middleware

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

//设置Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//参照：https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go

		//设置支持语言
		en := en.New()
		zh := zh.New()

		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		val := validator.New()

		//根据参数取翻译器实例
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		//翻译器注册到validator
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			break
		default:
			zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			//自定义验证方法
			//https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go
			val.RegisterValidation("is-validuser", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "admin"
			})

			//自定义验证器
			//https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
			val.RegisterTranslation("is-validuser", trans, func(ut ut.Translator) error {
				return ut.Add("is-validuser", "{0} 填写不正确哦", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("is-validuser", fe.Field())
				return t
			})

			// 验证服务名称格式
			val.RegisterValidation("valid_service_name", func(fl validator.FieldLevel) bool {
				isMatch, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
				return isMatch
			})
			val.RegisterTranslation("valid_service_name", trans, func(ut ut.Translator) error {
				return ut.Add("valid_service_name", "{0} 填写不符合格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_service_name", fe.Field())
				return t
			})

			// 验证接入规则格式
			val.RegisterValidation("valid_rule", func(fl validator.FieldLevel) bool {
				isMatch, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
				return isMatch
			})
			val.RegisterTranslation("valid_rule", trans, func(ut ut.Translator) error {
				return ut.Add("valid_rule", "{0} 填写不符合格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_rule", fe.Field())
				return t
			})

			// 验证url重写格式
			val.RegisterValidation("valid_url_rewrite", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, pair := range strings.Split(fl.Field().String(), ",") {
					if len(strings.Split(pair, " ")) != 2 {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("valid_url_rewrite", trans, func(ut ut.Translator) error {
				return ut.Add("valid_url_rewrite", "{0} 填写不符合格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_url_rewrite", fe.Field())
				return t
			})

			// 验证header转换格式
			val.RegisterValidation("valid_header_transform", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, pair := range strings.Split(fl.Field().String(), ",") {
					if len(strings.Split(pair, " ")) != 3 {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("valid_header_transform", trans, func(ut ut.Translator) error {
				return ut.Add("valid_header_transform", "{0} 填写不符合格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_header_transform", fe.Field())
				return t
			})

			// 验证ip:port格式
			val.RegisterValidation("valid_ip_port_list", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, addr := range strings.Split(fl.Field().String(), ",") {
					if isMatch, _ := regexp.Match(`^\S+\:\d+$`, []byte(addr)); !isMatch {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("valid_ip_port_list", trans, func(ut ut.Translator) error {
				return ut.Add("valid_ip_port_list", "{0} 填写不符合格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_ip_port_list", fe.Field())
				return t
			})

			// 验证ip格式
			val.RegisterValidation("valid_ip_list", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, addr := range strings.Split(fl.Field().String(), ",") {
					if isMatch, _ := regexp.Match(`\S+`, []byte(addr)); !isMatch {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("valid_ip_list", trans, func(ut ut.Translator) error {
				return ut.Add("valid_ip_list", "{0} 填写不符合格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_ip_list", fe.Field())
				return t
			})

			// 验证权重列表格式
			val.RegisterValidation("valid_weight_list", func(fl validator.FieldLevel) bool {
				for _, weight := range strings.Split(fl.Field().String(), ",") {
					if isMatch, _ := regexp.Match(`^\d+$`, []byte(weight)); !isMatch {
						return false
					}
				}
				return true
			})
			val.RegisterTranslation("valid_weight_list", trans, func(ut ut.Translator) error {
				return ut.Add("valid_weight_list", "{0} 填写不符合格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_weight_list", fe.Field())
				return t
			})

			break
		}
		c.Set(constants.TranslatorKey, trans)
		c.Set(constants.ValidatorKey, val)
		c.Next()
	}
}
