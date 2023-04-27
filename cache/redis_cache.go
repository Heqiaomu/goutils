package cache

import (
	"context"
	"fmt"
	log "github.com/Heqiaomu/glog"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type RedisTTLCache struct {
	Addr             string
	Name             string
	ttl              time.Duration // default 30s, for key expire time， if set zero and the key will not expire
	skipTTLExtension bool
	addrRefresh      func() string

	rdbc redis.Cmdable
}

type RedisTTLCacheOption func(*RedisTTLCache)

var defaultTTL = time.Second * 30

func NewRedisCache(addr, name string, opts ...RedisTTLCacheOption) (*RedisTTLCache, error) {
	return newRedisTTLCache(context.Background(), addr, name, defaultTTL, opts...)
}

func MustNewRedisCache(addr, name string, opts ...RedisTTLCacheOption) *RedisTTLCache {
	cache, err := newRedisTTLCache(context.Background(), addr, name, defaultTTL, opts...)
	if err != nil {
		panic(fmt.Sprintf("new redis client failed, Err: [%v]", err))
	}
	return cache
}

func WithSkipTTLExtensionOnHit() RedisTTLCacheOption {
	return func(cache *RedisTTLCache) {
		cache.skipTTLExtension = true
	}
}

func WithTTL(ttl time.Duration) RedisTTLCacheOption {
	return func(cache *RedisTTLCache) {
		cache.ttl = ttl
	}
}

func WithAddrRefresh(addrRefresh func() string) RedisTTLCacheOption {
	return func(cache *RedisTTLCache) {
		cache.addrRefresh = addrRefresh
	}
}

func newRedisTTLCache(ctx context.Context, addr, name string, ttl time.Duration, opts ...RedisTTLCacheOption) (*RedisTTLCache, error) {
	if name == "" {
		return nil, errors.New("fail to new redis cache, because name is empty, please check")
	}
	rtc := &RedisTTLCache{
		Addr: addr,
		Name: name,
		ttl:  ttl,
	}
	for _, opt := range opts {
		opt(rtc)
	}
	// if set ttl to zero, and the ttl extension is invalid, so assign true to skipTTLExtension
	if rtc.ttl == time.Duration(0) {
		rtc.skipTTLExtension = true
	}
	err := rtc.newClient(ctx, addr)
	if err != nil {
		return nil, errors.Wrap(err, "new client")
	}
	if rtc.addrRefresh != nil {
		go rtc.check(ctx)
	}
	return rtc, nil
}

func (rtc *RedisTTLCache) newClient(ctx context.Context, addr string) error {
	if addr == "" {
		return errors.New("fail to new redis cache, because addr is empty, please check")
	}
	var client redis.Cmdable
	if strings.Contains(addr, ",") || strings.Contains(addr, "redis-cluster") {
		clusters := strings.Split(addr, ",")
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    clusters,
			Password: "",
		}).WithContext(ctx)
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}).WithContext(ctx)
	}
	err := client.Ping(ctx).Err()
	if err != nil {
		return errors.Wrap(err, "ping")
	}
	rtc.rdbc = client
	return nil
}

func (rtc *RedisTTLCache) check(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			err := rtc.rdbc.Ping(ctx).Err()
			if err == nil {
				continue
			}
			err = rtc.newClient(ctx, rtc.addrRefresh())
			if err != nil {
				log.Warnf("Currently new redis client failed, Err: [%v]", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (rtc *RedisTTLCache) Close() {
	if rdc, ok := rtc.rdbc.(*redis.Client); ok {
		_ = rdc.Close()
	}
	if rdc, ok := rtc.rdbc.(*redis.ClusterClient); ok {
		_ = rdc.Close()
	}
}

// CtxSet set with context
func (rtc *RedisTTLCache) CtxSet(ctx context.Context, key, value string) error {
	err := rtc.rdbc.Set(ctx, fmt.Sprintf("{%s}:%s", rtc.Name, key), value, rtc.ttl).Err()
	if err != nil {
		return errors.Wrap(err, "set")
	}
	return nil
}

// Set set key value to redis with string type, ttl will use for this key if ttl var set,
// if ttl == 0 be set means the key has no expiration time.
func (rtc *RedisTTLCache) Set(key, value string) error {
	return rtc.CtxSet(context.Background(), key, value)
}

func (rtc *RedisTTLCache) CtxGet(ctx context.Context, key string) (string, error) {
	var result string
	var err error
	if rtc.skipTTLExtension {
		result, err = rtc.rdbc.Get(ctx, fmt.Sprintf("{%s}:%s", rtc.Name, key)).Result()
	} else {
		result, err = rtc.rdbc.GetEx(ctx, fmt.Sprintf("{%s}:%s", rtc.Name, key), rtc.ttl).Result()
	}
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", errors.Wrap(err, "get")
	}
	return result, nil
}

// Get get value with key, if value not exist, return empty string and nil error
func (rtc *RedisTTLCache) Get(key string) (string, error) {
	return rtc.CtxGet(context.Background(), key)
}

func (rtc *RedisTTLCache) CtxDel(ctx context.Context, key string) error {
	err := rtc.rdbc.Del(ctx, fmt.Sprintf("{%s}:%s", rtc.Name, key)).Err()
	if err != nil {
		return errors.Wrap(err, "del")
	}
	return nil
}

func (rtc *RedisTTLCache) Del(key string) error {
	return rtc.CtxDel(context.Background(), key)
}
