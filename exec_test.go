package util

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
)

//func TestCMD(t *testing.T) {
//	Convey(`Test cmdls -al cmd`, t, func() {
//		Convey(`When cmd='ls -al' no timeout`, func() {
//			out, err := CMD(2, "ls -al")
//			So(err, ShouldBeNil)
//			So(out, ShouldNotBeBlank)
//		})
//
//		Convey(`When cmd='sleep 2s' timeout`, func() {
//			out, err := CMD(1, "sleep 2s")
//			So(err, ShouldNotBeNil)
//			So(out, ShouldBeBlank)
//		})
//	})
//}

func TestSetShellType(t *testing.T) {
	type args struct {
		st ShellType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"TestSetShellType_1",
			args{
				ShellTypeBash,
			},
		},
		{
			"TestSetShellType_2",
			args{
				ShellTypeShell,
			},
		},
		{
			"TestSetShellType_3",
			args{
				ShellTypeNone,
			},
		},
		{
			"TestSetShellType_4",
			args{
				4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetShellType(tt.args.st)
		})
	}
}

func Test_setShellType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Test_setShellType_1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setShellType()
		})
	}
}

func Test_setShellType2(t *testing.T) {
	PathExists := ApplyFunc(PathExists, func(p string) (bool, error) {
		if p == "/bin/bash" {
			return false, nil
		}
		return true, nil
	})
	defer PathExists.Reset()

	tests := []struct {
		name string
	}{
		{
			"Test_setShellType_2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setShellType()
		})
	}
}

func Test_setShellType3(t *testing.T) {
	PathExists := ApplyFunc(PathExists, func(p string) (bool, error) {
		if p == "/bin/bash" || p == "/bin/sh" {
			return false, nil
		}
		return true, nil
	})
	defer PathExists.Reset()
	tests := []struct {
		name string
	}{
		{
			"Test_setShellType_3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setShellType()
		})
	}
}

func TestCommand(t *testing.T) {
	defer func() {
		os.RemoveAll("abc")
	}()

	p, _ := os.Getwd()

	type args struct {
		shellMode bool
		timeout   int
		forceKill bool
		command   string
		arg       []string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantRes bool
	}{
		{
			"TestCommand_1",
			args{
				true,
				10,
				true,
				"mkdir abc",
				[]string{},
			},
			"",
			true,
		},
		{
			"TestCommand_2",
			args{
				false,
				10,
				true,
				"pwd",
				[]string{},
			},
			p + "\n",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, gotRes := Command(tt.args.shellMode, tt.args.timeout, tt.args.forceKill, tt.args.command, tt.args.arg...)
			if gotOut != tt.wantOut {
				t.Errorf("Command() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
			if gotRes != tt.wantRes {
				t.Errorf("Command() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestNewDefaultCMD(t *testing.T) {
	mm := make(map[string]string)
	mm["key"] = "value"
	type args struct {
		command string
		args    []string
		opts    []CmdOption
	}
	tests := []struct {
		name string
		args args
		want *Ins
	}{
		{
			"TestNewDefaultCMD_1",
			args{
				"",
				[]string{"al"},
				[]CmdOption{
					WithRetry(5),
					WithEnvs(mm),
					WithShellType(ShellTypeShell),
					WithTimeout(5),
					WithForceKill(true),
					WithErrPrint(true),
					WithRetryInterval(time.Second),
				},
			},
			nil,
		},
		{
			"TestNewDefaultCMD_1",
			args{
				"ls",
				[]string{"al"},
				[]CmdOption{
					WithRetry(5),
					WithEnvs(mm),
					WithShellType(ShellTypeShell),
					WithTimeout(5),
					WithForceKill(true),
					WithErrPrint(true),
					WithRetry(2),
					WithRetryInterval(time.Second),
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := NewDefaultCMD(tt.args.command, tt.args.args, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewDefaultCMD() = %v, want %v", got, tt.want)
			// }
			got := NewDefaultCMD(tt.args.command, tt.args.args, tt.args.opts...)
			got.Run()
		})
	}
}

func TestExe(t *testing.T) {
	r, err := Exe("", "ls")
	assert.Nil(t, err)
	assert.NotNil(t, r)

	r, err = Exe("", "ls", "nonexistent_file")
	assert.NotNil(t, err)
	assert.Empty(t, r)
}
