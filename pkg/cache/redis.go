package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8"
)

type Config struct {
	// Host Port build Addr
	Host string
	// Port 端口
	Port string
	// Password 密码
	Password string `json:"password"`
	// DB，默认为0
	DB int `json:"db"`
	// PoolSize 最大连接池限制 默认每个CPU10个连接
	PoolSize int `json:"poolSize"`
	// MaxRetries 网络相关的错误最大重试次数 默认8次
	MaxRetries int `json:"maxRetries"`
	// MinIdleConns 最小空闲连接数
	MinIdleConns int `json:"minIdleConns"`
	// DialTimeout 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout"`
	// ReadTimeout 读超时 默认3s
	ReadTimeout time.Duration `json:"readTimeout"`
	// WriteTimeout 读超时 默认3s
	WriteTimeout time.Duration `json:"writeTimeout"`
	// IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	IdleTimeout time.Duration `json:"idleTimeout"`
	// Debug开关
	Debug bool `json:"debug"`
}

// redis distributed lock lua script
const script = `
	if redis.call('get', KEYS[1]) == ARGV[1] 
    then 
	    return redis.call('del', KEYS[1]) 
	else 
	    return 0 
	end
`

type Redis struct {
	*redis.Client
}

func NewRedis(cf *Config) (*Redis, error) {
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cf.Host, cf.Port),
		Password:     cf.Password,
		DB:           cf.DB,
		MaxRetries:   cf.MaxRetries,
		DialTimeout:  cf.DialTimeout,
		ReadTimeout:  cf.ReadTimeout,
		WriteTimeout: cf.WriteTimeout,
		PoolSize:     cf.PoolSize,
		MinIdleConns: cf.MinIdleConns,
		IdleTimeout:  cf.IdleTimeout,
	}
	if cf.MaxRetries > 0 {
		options.MaxRetries = cf.MaxRetries
	}
	if cf.DialTimeout > 0 {
		options.DialTimeout = cf.DialTimeout
	}
	if cf.ReadTimeout > 0 {
		options.ReadTimeout = cf.ReadTimeout
	}
	if cf.WriteTimeout > 0 {
		options.WriteTimeout = cf.WriteTimeout
	}
	if cf.PoolSize > 0 {
		options.PoolSize = cf.PoolSize
	}
	if cf.MinIdleConns > 0 {
		options.MinIdleConns = cf.MinIdleConns
	}
	if cf.IdleTimeout > 0 {
		options.IdleTimeout = cf.IdleTimeout
	}
	client := redis.NewClient(options)
	client.AddHook(apmgoredis.NewHook())
	_, err := client.Ping(context.Background()).Result()
	return &Redis{client}, err
}

// GetRedisLock redis 加锁
func (client *Redis) GetRedisLock(ctx context.Context, key string, expireTime int64, values ...interface{}) (interface{}, bool) {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	} else {
		u1, _ := uuid.NewUUID()
		value = u1.String()
	}
	result, _ := client.SetNX(ctx, key, value, time.Duration(expireTime)*time.Second).Result()
	if result {
		return value, true
	}
	return value, false
}

// ReleaseRedisLock redis 释放锁
func (client *Redis) ReleaseRedisLock(ctx context.Context, key string, value interface{}) bool {
	if result, err := client.Eval(ctx, script, []string{key}, value).Result(); err != nil {
		return false
	} else {
		val, ok := result.(int64)
		if !ok {
			return false
		}
		if val == 1 {
			return true
		}
		return false
	}

}
