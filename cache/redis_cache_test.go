package cache

import (
	"context"
	"fmt"
	log "github.com/Heqiaomu/glog"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestNewRedisCache(t *testing.T) {
	mock := func(t *testing.T) *gomonkey.Patches {
		patches := gomonkey.NewPatches()
		var c *redis.Client
		patches.ApplyMethod(reflect.TypeOf(c), "Process", func(c *redis.Client, ctx context.Context, cmd redis.Cmder) error {
			fmt.Println("af")
			return nil
		})

		patches.ApplyMethod(reflect.TypeOf(c), "WithContext", func(c *redis.Client, ctx context.Context) *redis.Client {
			fmt.Println("a12f")
			return c
		})
		patches.ApplyMethod(reflect.TypeOf(c), "Ping", func(c *redis.Client, ctx context.Context) *redis.StatusCmd {
			return &redis.StatusCmd{}
		})

		patches.ApplyMethod(reflect.TypeOf(c), "Close", func(c *redis.Client) error { return nil })

		var cl *redis.ClusterClient
		patches.ApplyMethod(reflect.TypeOf(cl), "Process", func(c *redis.ClusterClient, ctx context.Context, cmd redis.Cmder) error {
			return nil
		})
		patches.ApplyMethod(reflect.TypeOf(c), "Set", func(c *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			return &redis.StatusCmd{}
		})
		patches.ApplyMethod(reflect.TypeOf(c), "Get", func(c *redis.Client, ctx context.Context, key string) *redis.StringCmd {
			return &redis.StringCmd{}
		})
		patches.ApplyMethod(reflect.TypeOf(c), "GetEx", func(c *redis.Client, ctx context.Context, key string, expiration time.Duration) *redis.StringCmd {
			return &redis.StringCmd{}
		})

		patches.ApplyMethod(reflect.TypeOf(c), "Del", func(c *redis.Client, ctx context.Context, keys ...string) *redis.IntCmd {
			return &redis.IntCmd{}
		})
		patches.ApplyFunc(redis.NewClusterClient, func(opt *redis.ClusterOptions) *redis.ClusterClient {
			return &redis.ClusterClient{}
		})
		patches.ApplyMethod(reflect.TypeOf(cl), "Close", func(c *redis.ClusterClient) error { return nil })
		return patches
	}
	log.Logger()
	t.Run("test NewRedisCache", func(t *testing.T) {
		mk := mock(t)
		defer mk.Reset()
		cache, err := NewRedisCache("mockaddr", "mockname", WithTTL(0), WithSkipTTLExtensionOnHit())
		assert.Nil(t, err)
		assert.NotNil(t, cache)
		cache.Close()
		redisCache := MustNewRedisCache("127.0.0.1:6379", "mockname", WithTTL(0), WithSkipTTLExtensionOnHit())
		assert.NotNil(t, cache)
		//redisCache, err := NewRedisCache("redis-cluster", "mockname")
		//assert.Nil(t, err)
		//assert.NotNil(t, redisCache)

		err = redisCache.Set("a", "1")
		assert.Nil(t, err)
		get, err := redisCache.Get("a")
		assert.Nil(t, err)
		assert.NotNil(t, get)
		err = redisCache.Del("a")
		assert.Nil(t, err)
	})
}
