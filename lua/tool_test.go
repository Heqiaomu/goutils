package lua

import (
	"errors"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func TestMap2Table(t *testing.T) {
	mm := make(map[string]interface{})
	m2 := make(map[string]interface{})
	m3 := make(map[string]string)
	m4 := make([]map[string]interface{}, 1)
	m5 := make([]interface{}, 1)

	mm["k1"] = 1
	mm["k2"] = "2"
	mm["k3"] = true
	mm["k4"] = float64(1.0)
	mm["k5"] = []byte("a")
	mm["k6"] = m2
	mm["k7"] = m3
	mm["k8"] = time.Now()
	mm["k9"] = m4
	mm["k10"] = m5

	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *lua.LTable
	}{
		{
			"TestMap2Table_1",
			args{
				mm,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Map2Table(tt.args.m)
			// if got := Map2Table(tt.args.m); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Map2Table() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestErrorReturn(t *testing.T) {
	type args struct {
		l   *lua.LState
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"TestErrorReturn",
			args{
				lua.NewState(),
				errors.New("test error"),
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorReturn(tt.args.l, tt.args.err); got != tt.want {
				t.Errorf("ErrorReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataReturn(t *testing.T) {

	type arg struct {
		Key string
	}

	s := arg{"this is test case"}
	sArr := []arg{s}

	type args struct {
		l    *lua.LState
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"TestDataReturn_1",
			args{
				lua.NewState(),
				s,
			},
			2,
		},
		{
			"TestDataReturn_2",
			args{
				lua.NewState(),
				sArr,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DataReturn(tt.args.l, tt.args.data); got != tt.want {
				t.Errorf("DataReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}
