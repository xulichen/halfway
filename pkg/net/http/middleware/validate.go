package middleware

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	zhTranslations "github.com/go-playground/validator/translations/zh"
	"github.com/labstack/echo/v4"
	"github.com/xulichen/halfway/pkg/utils"
	"gopkg.in/go-playground/validator.v9"
)

// CustomBinder 自定bind
type CustomBinder struct{}

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
	e.Validator = utils.NewValidator()
	e.Binder = new(CustomBinder)
}
