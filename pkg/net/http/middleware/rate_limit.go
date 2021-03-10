package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"github.com/xulichen/halfway/pkg/consts"
	"github.com/xulichen/halfway/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

type LimitStoreOptionsConfig struct {
	Prefix string
	Client redis.Client
}

var limitStoreOptionsConfig LimitStoreOptionsConfig

func SetLimitStoreOptionsConfig(c LimitStoreOptionsConfig) {
	limitStoreOptionsConfig = c
}

// NewDefaultStore is in-memory store
func NewDefaultStore() limiter.Store {
	return memory.NewStore()
}

// newDefaultStoreOptions is 创建限流器默认配置
func newDefaultStoreOptions() *limiter.StoreOptions {
	storeOptions := new(limiter.StoreOptions)
	storeOptions.Prefix = limitStoreOptionsConfig.Prefix
	return storeOptions
}

// newApiRedisStoreLimiter is 基于redis的集群api限流器
func newApiRedisStoreLimiter(limit int64, period time.Duration) *limiter.Limiter {
	return newApiStoreLimiterWithOptions(limitStoreOptionsConfig.Client, newDefaultStoreOptions(), limit, period)
}

// newApiDefaultStoreLimiter is 基于单个服务的api限流器
func newApiDefaultStoreLimiter(limit int64, period time.Duration) *limiter.Limiter {
	return newApiStoreLimiterWithOptions(nil, newDefaultStoreOptions(), limit, period)
}

// NewApiRedisStoreLimiterWithOptions is 创建接口限流器
func newApiStoreLimiterWithOptions(
	client redis.Client, storeOptions *limiter.StoreOptions, limit int64, period time.Duration,
) *limiter.Limiter {
	var store limiter.Store
	if client != nil {
		redisStore, err := redis.NewStoreWithOptions(client, *storeOptions)
		if err != nil {
			panic(fmt.Sprintf("NewApiStoreLimiterWithOptions is failed, err is  %s", err.Error()))
		}
		store = redisStore
	} else {
		store = NewDefaultStore()
	}
	return limiter.New(store, limiter.Rate{
		Limit:  limit,
		Period: period,
	})
}

// api接口全局限流器
func RouteReachedLimitGlobal(limit int64, period time.Duration) echo.MiddlewareFunc {
	apiLimiter := newApiRedisStoreLimiter(limit, period)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// 接口是否达到限流
			ctx, err := apiLimiter.Get(c.Request().Context(), fmt.Sprintf("%s_%s", c.Request().URL.Path, c.Request().Method))
			if err != nil {
				utils.GetLogger().Errorf("限流器异常 - err: %v, %s on %s", err, ctx, c.Request().URL)
				return next(c)
			}
			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(ctx.Limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(ctx.Remaining, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(ctx.Reset, 10))
			if ctx.Reached { // 钉钉调用已到上限
				utils.GetLogger().Errorf(" %s is failed reached limit. limit is %d, remaining is %d, reset is %d", c.Request().URL, ctx.Limit, ctx.Remaining, ctx.Reset)
				return c.JSON(http.StatusTooManyRequests, echo.Map{
					"code": consts.ResponseCodeErrParameter,
					"msg":  consts.ResponseRateLimitReachedStatusText,
				})
			}

			return next(c)
		}
	}
}

// api接口本地限流器
func RouteReachedLimitLocal(limit int64, period time.Duration) echo.MiddlewareFunc {
	apiLimiter := newApiDefaultStoreLimiter(limit, period)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// 接口是否达到限流
			ctx, err := apiLimiter.Get(c.Request().Context(), fmt.Sprintf("%s_%s", c.Request().URL.Path, c.Request().Method))
			if err != nil {
				utils.GetLogger().Errorf("限流器异常 - err: %v, %s on %s", err, ctx, c.Request().URL)
				return next(c)
			}
			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(ctx.Limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(ctx.Remaining, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(ctx.Reset, 10))
			if ctx.Reached { // 钉钉调用已到上限
				return c.JSON(http.StatusTooManyRequests, echo.Map{
					"code": consts.ResponseCodeErrParameter,
					"msg":  consts.ResponseRateLimitReachedStatusText,
				})
			}

			return next(c)
		}
	}
}

// RouteUserLimitGlobal api接口单用户全局限流器
func RouteUserLimitGlobal(limit int64, period time.Duration) echo.MiddlewareFunc {
	apiLimiter := newApiRedisStoreLimiter(limit, period)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// 接口是否达到限流
			ctx, err := apiLimiter.Get(c.Request().Context(), c.Request().Header.Get("Authorization"))
			if err != nil {
				utils.GetLogger().Errorf("限流器异常 - err: %v, %s on %s", err, ctx, c.Request().URL)
				return next(c)
			}
			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(ctx.Limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(ctx.Remaining, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(ctx.Reset, 10))
			if ctx.Reached {
				return c.JSON(http.StatusOK, echo.Map{
					"code": consts.ResponseCodeErrParameter,
					"msg":  consts.ResponseRateLimitReachedStatusText,
				})
			}

			return next(c)
		}
	}
}
