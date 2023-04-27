package process

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Running 判断是否正在运行
func Running(bin string) (run bool, pid string) {
	output, _ := exec.Command("pgrep", "-f", bin).Output()
	pidStr := strings.TrimSpace(string(output))

	return !(pidStr == ""), pidStr
}

// RunClearlyWithPWD with pwd
func RunClearlyWithPWD(bin string, args []string, dir string, std ...io.Writer) (cmd *exec.Cmd, err error) {
	cmd = exec.Command(bin, args...)
	if strings.TrimSpace(dir) != "" {
		cmd.Dir = dir
	}

	if err = addLDLibraryPathToCMDEnv(bin, cmd); err != nil {
		return
	}

	if len(std) == 1 {
		cmd.Stdout = std[0]
		cmd.Stderr = std[0]
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err = cmd.Start(); err != nil {
		return
	}

	if s, _ := started(bin); !s {
		err = fmt.Errorf("fail to start executable [%s], because timeout", bin)
		return
	}
	pid := cmd.Process.Pid
	fmt.Printf("RunClearly. pid=%d\n", pid)
	return
}

// RunClearly 启动app没有僵尸进程
func RunClearly(bin string, args []string, std ...io.Writer) (cmd *exec.Cmd, err error) {
	cmd = exec.Command(bin, args...)

	if err = addLDLibraryPathToCMDEnv(bin, cmd); err != nil {
		return
	}

	if len(std) == 1 {
		cmd.Stdout = std[0]
		cmd.Stderr = std[0]
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err = cmd.Start(); err != nil {
		return
	}

	if s, _ := started(bin); !s {
		err = fmt.Errorf("start executable timeout: %s", bin)
		return
	}
	pid := cmd.Process.Pid
	fmt.Printf("RunClearly. pid=%d\n", pid)
	return
}

func addLDLibraryPathToCMDEnv(bin string, cmd *exec.Cmd) error {
	path := bin[:strings.LastIndex(bin, "/")]
	env := os.Environ()
	var cmdEnv []string
	var exist bool
	for _, e := range env {
		i := strings.Index(e, "=")
		if i > 0 && (e[:i] == "LD_LIBRARY_PATH") {
			exist = true
			cmdEnv = append(cmdEnv, fmt.Sprintf("LD_LIBRARY_PATH=%s:%s", e[i+1:], path))
		} else {
			cmdEnv = append(cmdEnv, e)
		}
	}
	if !exist {
		cmdEnv = append(cmdEnv, fmt.Sprintf("LD_LIBRARY_PATH=%s", path))
	}
	level := viper.GetString("log.console_level")
	if len(level) != 0 {
		cmdEnv = append(cmdEnv, fmt.Sprintf("CONSOLE_LOGGER_LEVEL=%s", level))
	}
	cmd.Env = cmdEnv
	return nil
}

// Run 启动app
func Run(bin string, args []string) (pid int, err error) {
	cmd := exec.Command(bin, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		return 0, err
	}

	if s, _ := started(bin); !s {
		return 0, fmt.Errorf("start executable timeout: %s", bin)
	}
	pid = cmd.Process.Pid

	return
}

// 判断是否已经启动成功
func started(bin string) (s bool, pid string) {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if r, pid := Running(bin); r {
				return true, pid
			}
		case <-time.After(time.Second * 2):
			return false, ""
		}
	}
}

// Stop app
func Stop(bin string) (err error) {
	run, pid := Running(bin)
	if !run {
		return nil
	}

	cmd := exec.Command("kill", "-TERM", pid)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func StopPid(pid int) (err error) {
	if pid == 0 {
		return fmt.Errorf("pid is 0")
	}
	cmd := exec.Command("kill", "-TERM", strconv.Itoa(pid))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("StopPid. fail to run 'cmd.Run()'. err=%v\n", err)
		return
	}
	return
}
