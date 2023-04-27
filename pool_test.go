package util

import (
	"testing"
)

func TestNewPool(t *testing.T) {
	type args struct {
		factory Factory
		destroy Destory
		cap     int
	}
	tests := []struct {
		name string
		args args
		want *Pool
	}{
		{
			"TestNewPool_1",
			args{
				func() (ConnRes, error) {
					return "", nil
				},
				func(ConnRes) {
					return
				},
				1,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPool(tt.args.factory, tt.args.destroy, tt.args.cap)
			got.Put("this is one")
			got.Put("this is two")
			got.Get()
			got.Get()
			got.Get()
			got.DestoryPool()
		})
	}
}
