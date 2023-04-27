package distributelock

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/client/pkg/v3/logutil"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"regexp"
	"sync"
	"time"
)

var dc *clientv3.Client
var mu sync.Mutex

// DistributeLock .
type DistributeLock struct {
	//c        *clientv3.Client // client for etcd
	etcdURLs []string      // etcd urls
	lockKey  string        // lockKey for lock in etcd
	timeout  time.Duration // timout for Lock, Use WithTimeout to set this
	ctx      context.Context

	//Release // Release distributed lock manager, don't forget to call this
}

// DistributeLockOpt .
type DistributeLockOpt func(lock *DistributeLock)

// WithTimeout timeout option
func WithTimeout(duration time.Duration) func(lock *DistributeLock) {
	return func(lock *DistributeLock) {
		lock.timeout = duration
	}
}

// WithEtcdURLs etcd urls option
func WithEtcdURLs(urls []string) func(lock *DistributeLock) {
	return func(lock *DistributeLock) {
		lock.etcdURLs = urls
	}
}

type UnLock func() error

var keyreg = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9\\.\\-_]{3,63}[a-zA-Z0-9]$")

// NewDistributeLock new distribute lock
// key in format of '^[a-zA-Z][a-zA-Z0-9\.\-_]{3,63}[a-zA-Z0-9]$'
// full key is '/distributelock/{key}/{leaseid}'
func NewDistributeLock(key string, opts ...DistributeLockOpt) (*DistributeLock, error) {
	if !keyreg.MatchString(key) {
		return nil, errors.Errorf(fmt.Sprintf("fail to validate key [%s] format, bacause not match regex '^[a-zA-Z][a-zA-Z0-9\\.\\-_]{3,63}[a-zA-Z0-9]$', please input correct key", key))
	}
	fkey := fmt.Sprintf("/distributelock/%s", key)
	dl := &DistributeLock{ctx: context.Background(), etcdURLs: viper.GetStringSlice("etcd.host"), lockKey: fkey}
	for _, opt := range opts {
		opt(dl)
	}
	err := checkdc(dl)
	if err != nil {
		return nil, errors.Wrap(err, "init etcd cli")
	}
	// 每次都新建链接这个肯定不行
	//cli, err := clientv3.New(config)
	//if err != nil {
	//	return nil, errors.Wrap(err, "new etcd client")
	//}
	//dl.c = cli
	//dl.Release = func() error { return cli.Close() }
	return dl, nil
}

func checkdc(dl *DistributeLock) error {
	if dc == nil {
		mu.Lock()
		defer mu.Unlock()
		if dc == nil {
			loggerConfig := logutil.DefaultZapLoggerConfig
			loggerConfig.Sampling = nil
			config := clientv3.Config{
				Context:     dl.ctx,
				Endpoints:   dl.etcdURLs,
				DialTimeout: time.Second * 3,
				LogConfig:   &loggerConfig,
			}
			cli, err := clientv3.New(config)
			if err != nil {
				return errors.Wrap(err, "new etcd client")
			}
			dc = cli
		}
	} else {
		mu.Lock()
		defer mu.Unlock()
		if dc != nil {
			timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
			defer cancelFunc()
			if err := dc.Sync(timeout); err != nil {
				_ = dc.Close()
				dc = nil
				return errors.Wrap(err, "check etcd connection")
			}
		} else {
			return errors.New("cli is nil")
		}
	}
	return nil
}

// Lock lock with etcd, return unLockFunc and error
// dont forget to call unLockFunc when you finish
// your work
func (d *DistributeLock) Lock() (unLockFunc UnLock, done <-chan struct{}, err error) {
	err = checkdc(d)
	if err != nil {
		return nil, nil, errors.Wrap(err, "check etcd cli")
	}
	session, err := concurrency.NewSession(dc, concurrency.WithTTL(10))
	if err != nil {
		return nil, nil, errors.Wrap(err, "new lock session, please check the lock is released in advance")
	}
	mutex := concurrency.NewMutex(session, d.lockKey)
	ctx := d.ctx
	if d.timeout != time.Duration(0) {
		tctx, cancelFunc := context.WithTimeout(d.ctx, d.timeout)
		defer cancelFunc()
		ctx = tctx
	}
	err = mutex.Lock(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "lock")
	}
	return func() error {
		timeout, cancelFunc := context.WithTimeout(d.ctx, time.Second*3)
		defer cancelFunc()
		defer session.Close()
		return mutex.Unlock(timeout)
	}, session.Done(), nil
}
