package middleware

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	zhTranslations "github.com/go-playground/validator/translations/zh"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

type (
	// CustomBinder 自定bind
	CustomBinder struct{}

	// CustomValidator 自定义验证器结构体
	CustomValidator struct {
		validate *validator.Validate
		trans    ut.Translator
	}
)

// 自定义验证器
func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.validate.Struct(i); err != nil {
		var errList []string
		for _, e := range err.(validator.ValidationErrors) {
			errList = append(errList, e.Translate(c.trans))
		}
		return fmt.Errorf("%s", strings.Join(errList, "|"))
	}
	return nil
}

// Bind 自定义bind
func (cb *CustomBinder) Bind(i interface{}, e echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, e); err != nil {
		return fmt.Errorf("input error")
	}
	if err := e.Validate(i); err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	return
}

func InitValidate(e *echo.Echo) {
	zhCh := zh.New()
	uni := ut.New(zhCh)
	trans, _ := uni.GetTranslator("zh_Hans_CN")
	validate := validator.New()
	if err := zhTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic("注册本地化验证错误失败")
	}
	e.Validator = &CustomValidator{validate: validate, trans: trans}
	e.Binder = new(CustomBinder)
}
