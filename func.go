package util

import (
	"runtime"
	"strings"
)

// RunFuncName 运行函数名的全路径
func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// RunFuncName 运行函数的名字
func RunFuncNameEx() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return strings.Split(f.Name(), ".")[(len(strings.Split(f.Name(), ".")) - 1)]
}
