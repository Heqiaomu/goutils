package util

import (
	"testing"
)

func TestRunFuncName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"TestRunFuncName_1",
			"github.com/Heqiaomu/goutil.TestRunFuncName.func1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunFuncName(); got != tt.want {
				t.Errorf("RunFuncName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunFuncNameEx(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"TestRunFuncNameEx_1",
			"func1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunFuncNameEx(); got != tt.want {
				t.Errorf("RunFuncNameEx() = %v, want %v", got, tt.want)
			}
		})
	}
}
