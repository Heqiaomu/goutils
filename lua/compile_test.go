package lua

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestCompileFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *lua.FunctionProto
		wantErr bool
	}{
		{
			"TestCompileFile",
			args{
				"compile.go",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CompileFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompileFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCompileString(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    *lua.FunctionProto
		wantErr bool
	}{
		{
			"TestCompileString",
			args{
				"compile.go",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CompileString(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompileString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDoCompiledFile(t *testing.T) {
	type args struct {
		L     *lua.LState
		proto *lua.FunctionProto
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DoCompiledFile(tt.args.L, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("DoCompiledFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
