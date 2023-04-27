package distributelock

import (
	"context"
	log "github.com/Heqiaomu/glog"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"reflect"
	"testing"
	"time"
)

func TestNewDistributeLock(t *testing.T) {
	mock := func(t *testing.T) *gomonkey.Patches {
		patches := gomonkey.NewPatches()
		patches.ApplyFunc(clientv3.New, func(cfg clientv3.Config) (*clientv3.Client, error) {
			return &clientv3.Client{}, nil
		})
		var c *clientv3.Client
		patches.ApplyMethod(reflect.TypeOf(c), "Sync", func(c *clientv3.Client, ctx context.Context) error {
			return nil
		})
		return patches
	}
	log.Logger()
	t.Run("test NewDistributeLock", func(t *testing.T) {
		mk := mock(t)
		defer mk.Reset()
		lock, err := NewDistributeLock("KK", WithTimeout(time.Minute), WithEtcdURLs([]string{"127.0.0.1:2379"}))
		assert.NotNil(t, err)
		assert.Nil(t, lock)

		lock, err = NewDistributeLock("KKdadad", WithTimeout(time.Minute), WithEtcdURLs([]string{"127.0.0.1:2379"}))
		assert.Nil(t, err)
		assert.NotNil(t, lock)
	})
}

func TestLock(t *testing.T) {
	mock := func(t *testing.T) *gomonkey.Patches {
		patches := gomonkey.NewPatches()
		patches.ApplyFunc(concurrency.NewSession, func(client *clientv3.Client, opts ...concurrency.SessionOption) (*concurrency.Session, error) {
			return &concurrency.Session{}, nil
		})
		patches.ApplyFunc(concurrency.NewMutex, func(s *concurrency.Session, pfx string) *concurrency.Mutex {
			return &concurrency.Mutex{}
		})
		var m *concurrency.Mutex
		patches.ApplyMethod(reflect.TypeOf(m), "Lock", func(m *concurrency.Mutex, ctx context.Context) error {
			return nil
		})
		patches.ApplyMethod(reflect.TypeOf(m), "Unlock", func(m *concurrency.Mutex, ctx context.Context) error {
			return nil
		})
		var c *clientv3.Client
		patches.ApplyMethod(reflect.TypeOf(c), "Sync", func(c *clientv3.Client, ctx context.Context) error {
			return nil
		})
		return patches
	}
	log.Logger()
	t.Run("test Lock", func(t *testing.T) {
		mk := mock(t)
		defer mk.Reset()
		lock := DistributeLock{ctx: context.Background(), timeout: time.Second}
		lock.Lock()
	})
}

//
//func TestNewDistributeLoc1k(t *testing.T) {
//	mock := func(t *testing.T) *gomonkey.Patches {
//		patches := gomonkey.NewPatches()
//
//		return patches
//	}
//	log.Logger()
//	t.Run("test NewDistributeLock", func(t *testing.T) {
//		mk := mock(t)
//		defer mk.Reset()
//		lock, err := NewDistributeLock("test1_1", WithTimeout(time.Second*20), WithEtcdURLs([]string{"127.0.0.1:2379"}))
//		assert.Nil(t, err)
//		eg := sync.WaitGroup{}
//		for i := 0; i < 3; i++ {
//			eg.Add(1)
//			go func(i int) {
//				defer eg.Done()
//				unlock, done, e := lock.Lock()
//				if e != nil {
//					t.Logf("%d lock error, err : [%v]", i, e)
//					return
//				}
//				go func() {
//					select {
//					case <-done:
//						t.Logf("[%d] get session done info", i)
//					}
//				}()
//
//				t.Logf("[%d] got lock", i)
//				time.Sleep(time.Duration(2 * time.Second))
//
//				e = unlock()
//				if e != nil {
//					t.Logf("%d unlock error, err : [%v]", i, e)
//					return
//				}
//				t.Logf("[%d] unlock", i)
//			}(i)
//			//time.Sleep(time.Millisecond * 234)
//		}
//		for i := 3; i < 6; i++ {
//			eg.Add(1)
//			go func(i int) {
//				defer eg.Done()
//				unlock, done, e := lock.Lock()
//				if e != nil {
//					t.Logf("%d lock error, err : [%v]", i, e)
//					return
//				}
//				go func() {
//					select {
//					case <-done:
//						t.Logf("[%d] get session done info", i)
//
//					}
//				}()
//
//				t.Logf("[%d] got lock", i)
//				time.Sleep(time.Duration(2 * time.Second))
//				e = unlock()
//				if e != nil {
//					t.Logf("%d unlock error, err : [%v]", i, e)
//					return
//				}
//				t.Logf("[%d] unlock", i)
//			}(i)
//			//time.Sleep(time.Millisecond * 234)
//		}
//		eg.Wait()
//		//lock.Release()
//	})
//}
