package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/xulichen/halfway/pkg/utils"
	"runtime"
)

type (
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echoMiddleware.Skipper

		// Size of the stack to be printed.
		// Optional. Default value 4KB.
		StackSize int `yaml:"stack_size"`

		// DisableStackAll disables formatting stack traces of all other goroutines
		// into buffer after the trace for the current goroutine.
		// Optional. Default value false.
		DisableStackAll bool `yaml:"disable_stack_all"`

		// DisablePrintStack disables printing stack trace.
		// Optional. Default value as false.
		DisablePrintStack bool `yaml:"disable_print_stack"`

		// LogLevel is log level to printing stack trace.
		// Optional. Default value 0 (Print).
		LogLevel log.Lvl
	}
)

var (
	// DefaultRecoverConfig is the default Recover middleware config.
	defaultRecoverConfig = RecoverConfig{
		Skipper:           echoMiddleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          4,
	}
	errNotifyUrl string
)

// 报错通知地址 不设置不发送
func SetErrNotifyUrl(s string) {
	errNotifyUrl = s
}

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover(config *RecoverConfig) echo.MiddlewareFunc {
	return RecoverWithConfig(config)
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config *RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config == nil {
		config = &defaultRecoverConfig
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
						_ = utils.SendDingMsgWithUrl(errNotifyUrl, msg, true)
						switch config.LogLevel {
						case log.DEBUG:
							utils.GetLogger().Debug(msg)
						case log.INFO:
							utils.GetLogger().Info(msg)
						case log.WARN:
							utils.GetLogger().Warn(msg)
						case log.ERROR:
							utils.GetLogger().Error(msg)
						case log.OFF:
							// None.
						default:
							utils.GetLogger().Print(msg)
						}
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
