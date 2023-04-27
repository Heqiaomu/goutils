package process

import (
	"os/exec"
	"testing"
)

func TestRunning(t *testing.T) {
	type args struct {
		bin string
	}
	tests := []struct {
		name    string
		args    args
		wantRun bool
		wantPid string
	}{
		{
			"TestRunning",
			args{
				"xxx",
			},
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRun, gotPid := Running(tt.args.bin)
			if gotRun != tt.wantRun {
				t.Errorf("Running() gotRun = %v, want %v", gotRun, tt.wantRun)
			}
			if gotPid != tt.wantPid {
				t.Errorf("Running() gotPid = %v, want %v", gotPid, tt.wantPid)
			}
		})
	}
}

func TestRunClearly(t *testing.T) {
	type args struct {
		bin  string
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantCmd *exec.Cmd
		wantErr bool
	}{
		{
			"TestRunClearly",
			args{
				"a/b/c/d",
				[]string{"arg1", "arg2"},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RunClearly(tt.args.bin, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunClearly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRunClearlyWithPwd(t *testing.T) {
	type args struct {
		bin  string
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantCmd *exec.Cmd
		wantErr bool
	}{
		{
			"TestRunClearlyWithPwd",
			args{
				"a/b/c/d",
				[]string{"arg1", "arg2"},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RunClearlyWithPWD(tt.args.bin, tt.args.args, "mock")
			if (err != nil) != tt.wantErr {
				t.Errorf("RunClearlyWithPwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		bin  string
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantPid int
		wantErr bool
	}{
		{
			"TestRun",
			args{
				"grep",
				[]string{"-rnl", "a", "../"},
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.args.bin, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			StopPid(got)
		})
	}
}

func TestStop(t *testing.T) {
	type args struct {
		bin string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestStop",
			args{
				"grep",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Stop(tt.args.bin); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
