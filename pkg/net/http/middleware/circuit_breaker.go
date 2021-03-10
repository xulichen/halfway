// 断路器

package middleware

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xulichen/halfway/pkg/consts"
	"github.com/xulichen/halfway/pkg/log"
	"github.com/xulichen/halfway/pkg/net/http/errors"
)

// CircuitBreakerConfig 断路器配置 Name 需保证唯一
type CircuitBreakerConfig struct {
	Skipper                middleware.Skipper
	Name                   string `json:"name"` // should be unique
	Timeout                int    `json:"timeout"`
	MaxConcurrentRequests  int    `json:"max_concurrent_requests"`
	RequestVolumeThreshold int    `json:"request_volume_threshold"`
	SleepWindow            int    `json:"sleep_window"`
	ErrorPercentThreshold  int    `json:"error_percent_threshold"`
}

// CircuitBreaker 断路器中间件
func CircuitBreaker(config CircuitBreakerConfig) echo.MiddlewareFunc {
	var commendConfig hystrix.CommandConfig
	if config.Timeout != 0 {
		commendConfig.Timeout = config.Timeout
	}
	if config.MaxConcurrentRequests != 0 {
		commendConfig.MaxConcurrentRequests = config.MaxConcurrentRequests
	}
	if config.RequestVolumeThreshold != 0 {
		commendConfig.RequestVolumeThreshold = config.RequestVolumeThreshold
	}
	if config.SleepWindow != 0 {
		commendConfig.SleepWindow = config.SleepWindow
	}
	if config.ErrorPercentThreshold != 0 {
		commendConfig.ErrorPercentThreshold = config.ErrorPercentThreshold
	}

	hystrix.ConfigureCommand(config.Name, commendConfig)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) (err error) {
			// 异步模式
			output := make(chan bool, 0)
			e := hystrix.Go(config.Name, func() error {
				defer func() {
					if err := recover(); err != nil {
						log.GetLogger().Error(err)
					}
				}()
				if err = next(context); err == errors.ErrOutsideServerError {
					return err
				}
				output <- true
				return nil
			}, nil)
			select {
			case err = <-e:
				if err == hystrix.ErrCircuitOpen {
					return context.JSON(200, echo.Map{"code": consts.ResponseCodeInternalServerError, "msg": "circuit open"})
				} else if err == hystrix.ErrMaxConcurrency {
					return context.JSON(200, echo.Map{"code": consts.ResponseCodeInternalServerError, "msg": "max concurrency"})
				}
				return err
			case <-output:
				return nil
			}

			// 同步模式
			//err = hystrix.Do(config.Name, func() error {
			//	if err = next(context); err == entity.ErrOutsideServerError {
			//		return err
			//	}
			//	return nil
			//}, nil)
			//switch err {
			//case hystrix.ErrCircuitOpen:
			//	return context.JSON(200, echo.Map{"code": consts.ResponseCodeInternalServerError, "msg": "circuit open"})
			//case hystrix.ErrMaxConcurrency:
			//	return context.JSON(200, echo.Map{"code": consts.ResponseCodeInternalServerError, "msg": "max concurrency"})
			//default:
			//	return nil
			//}
		}
	}
}
