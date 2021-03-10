package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	zhTranslations "github.com/go-playground/validator/translations/zh"
	"gopkg.in/go-playground/validator.v9"
)

// CustomValidator 自定义验证器结构体
type CustomValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidator() *CustomValidator {
	zhCh := zh.New()
	uni := ut.New(zhCh)
	trans, _ := uni.GetTranslator("zh_Hans_CN")
	validate := validator.New()
	if err := zhTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic("注册本地化验证错误失败")
	}
	return &CustomValidator{
		validate: validate,
		trans:    trans,
	}
}

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
