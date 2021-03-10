package errors

import "errors"

var (
	// ErrOutsideServerError 抛出外部服务异常
	ErrOutsideServerError = errors.New("外部服务异常")
	// ErrInternalServerError 抛出内部服务异常
	_ = errors.New("内部服务错误")
)
