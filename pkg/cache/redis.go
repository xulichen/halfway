package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8"
	"github.com/xulichen/halfway/pkg/config"
	"time"
)

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

func NewRedis(cf *config.RedisConfig, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", cf.Host, cf.Port),
		Password:   cf.Password,
		DB:         db,
		MaxRetries: 5,
	})
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