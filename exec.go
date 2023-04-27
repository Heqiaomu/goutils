//Package util util
package util

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

var shellType ShellType

// SetShellType 设置全局的Shell类型 在执行CMD时使用该全局设置的Shell类型
func SetShellType(st ShellType) {
	if st == ShellTypeBash {
		shellType = ShellTypeBash
	} else if st == ShellTypeShell {
		shellType = ShellTypeShell
	} else if st == ShellTypeNone {
		shellType = ShellTypeNone
	} else {
		setShellType()
	}
}

// CMD run command
// timeout 设置超时时间，单位 秒
// command shell脚本或二进制全路径
// arg 二进制运行参数，shell模式无效
func CMD(timeout int, command string, arg ...string) (out string, error error) {
	// 如果 shellType 还没有被设置，那么根据运行系统中是否存在bash、sh来选择shell类型
	if shellType == 0 {
		setShellType()
	}
	out, ok := Cmd(shellType, timeout, true, command, arg...)
	if !ok {
		return out, fmt.Errorf("fail to exe %s", command)
	}
	return out, nil
}

func setShellType() {
	if exist, _ := PathExists("/bin/bash"); exist {
		shellType = ShellTypeBash
	} else if exist, _ := PathExists("/bin/sh"); exist {
		shellType = ShellTypeShell
	} else {
		shellType = ShellTypeNone
	}
}

// Command run command
// shellMode 是否是shell脚本
// timeout 设置超时时间
// forceKill 达到超时时间，是否强制杀死进程（包括子进程和孙进程）
// command shell脚本或二进制全路径
// arg 二进制运行参数，shell模式无效
func Command(shellMode bool, timeout int, forceKill bool, command string, arg ...string) (out string, res bool) {
	if shellMode {
		return Cmd(ShellTypeBash, timeout, forceKill, command, arg...)
	} else {
		return Cmd(ShellTypeNone, timeout, forceKill, command, arg...)
	}
}

type ShellType int

const (
	ShellTypeNone  = 1
	ShellTypeShell = 2
	ShellTypeBash  = 3
)

func Cmd(shellType ShellType, timeout int, forceKill bool, command string, arg ...string) (out string, res bool) {
	var (
		b           bytes.Buffer
		cmd         *exec.Cmd
		sysProcAttr *syscall.SysProcAttr
	)

	sysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // 使子进程拥有自己的 pgid，等同于子进程的 pid
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	switch shellType {
	case ShellTypeNone:
		cmd = exec.Command(command, arg...)
	case ShellTypeShell:
		cmd = exec.Command("/bin/sh", "-c", command)
	case ShellTypeBash:
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	cmd.SysProcAttr = sysProcAttr
	cmd.Stdout = &b
	cmd.Stderr = &b

	if err := cmd.Start(); err != nil {
		fmt.Printf("%s\n%s\n", b.String(), err.Error())
		return
	}

	waitChan := make(chan struct{}, 1)
	defer close(waitChan)

	// 超时杀掉进程组 或正常退出
	go func() {
		select {
		case <-ctx.Done():
			//fmt.Println("ctx timeout")
			if forceKill {
				//fmt.Printf("ctx timeout kill job ppid:%d\n", cmd.Process.Pid)
				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
					//fmt.Println("syscall.Kill return err: ", err)
					return
				}
			}
		case <-waitChan:
			//fmt.Printf("normal quit job ppid:%d\n", cmd.Process.Pid)
		}
	}()

	if err := cmd.Wait(); err != nil {
		//fmt.Printf("timeout kill job ppid:%s\n%s\n", b.String(), err.Error())
		em := err.Error()
		// 超时退出，返回调用失败
		if strings.Contains(em, "signal: killed") {
			return "", false
		}
		// 未超时，被执行程序主动退出
		return fmt.Sprintf("%s\n%v", b.String(), err), true
	}

	out = b.String()
	res = true
	// 唤起正常退出
	waitChan <- struct{}{}

	return
}

//=======================  new command func ==========================

type Ins struct {
	shellType     ShellType
	envs          map[string]string
	timeout       int
	forceKill     bool
	command       string
	args          []string
	retry         int
	retryInterval time.Duration
	errPrint      bool
	workDir       string
}
type CmdOption func(info *Ins)

func NewDefaultCMD(command string, args []string, opts ...CmdOption) *Ins {
	setShellType()
	i := Ins{
		shellType:     shellType,
		timeout:       10,
		forceKill:     true,
		command:       command,
		args:          args,
		retry:         1,
		retryInterval: time.Second,
	}
	for _, opt := range opts {
		opt(&i)
	}
	return &i
}

func WithRetry(retry int) CmdOption {
	if retry < 0 {
		retry = 1
	}
	return func(i *Ins) {
		i.retry = retry
	}
}

func WithWorkDir(wd string) CmdOption {
	return func(i *Ins) {
		i.workDir = wd
	}
}

func WithEnvs(envs map[string]string) CmdOption {
	return func(i *Ins) {
		i.envs = envs
	}
}

func WithShellType(shellType ShellType) CmdOption {
	return func(i *Ins) {
		i.shellType = shellType
	}
}

func WithTimeout(timeout int) CmdOption {
	return func(i *Ins) {
		i.timeout = timeout
	}
}

func WithForceKill(forceKill bool) CmdOption {
	return func(i *Ins) {
		i.forceKill = forceKill
	}
}

func WithErrPrint(ep bool) CmdOption {
	return func(i *Ins) {
		i.errPrint = ep
	}
}

func WithRetryInterval(interval time.Duration) CmdOption {
	return func(i *Ins) {
		i.retryInterval = interval
	}
}

func (i *Ins) Run() (out string, err error) {
	for t := 0; t < i.retry; t++ {
		out, err = i.cmd()
		if err != nil {
			if i.errPrint {
				fmt.Printf("fail to exe : %s\n", i.command)
				fmt.Printf("output : %s\n", out)
				fmt.Printf("err : %s\n", err.Error())
				fmt.Printf("retry : %d\n", t+1)
			}
			if i.retry != 1 {
				time.Sleep(i.retryInterval)
			}
		} else {
			break
		}
	}
	return
}

func (i *Ins) cmd() (out string, err error) {
	var (
		stdout      bytes.Buffer
		stderr      bytes.Buffer
		cmd         *exec.Cmd
		sysProcAttr *syscall.SysProcAttr
	)

	sysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // 使子进程拥有自己的 pgid，等同于子进程的 pid
	}

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.timeout)*time.Second)
	defer cancel()

	switch i.shellType {
	case ShellTypeNone:
		cmd = exec.Command(i.command, i.args...)
	case ShellTypeShell:
		cmd = exec.Command("/bin/sh", "-c", i.command)
	case ShellTypeBash:
		cmd = exec.Command("/bin/bash", "-c", i.command)
	}
	if cmd == nil {
		return "", errors.New("fail to new command")
	}

	cmd.SysProcAttr = sysProcAttr
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if i.workDir != "" {
		cmd.Dir = i.workDir
	}

	// add env to cmd
	if i.envs != nil {
		addEnv(cmd, i.envs)
	}

	if err := cmd.Start(); err != nil {
		if i.errPrint {
			fmt.Printf("stdout: %s\n", stdout.String())
			fmt.Printf("stderr: %s\n", stderr.String())
			fmt.Printf("err: %s\n", err.Error())
		}
		return "", err
	}

	waitChan := make(chan struct{}, 1)
	defer close(waitChan)

	// 超时杀掉进程组 或正常退出
	go func() {
		select {
		case <-ctx.Done():
			//fmt.Println("ctx timeout")
			if i.forceKill {
				//fmt.Printf("ctx timeout kill job ppid:%d\n", cmd.Process.Pid)
				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
					//fmt.Println("syscall.Kill return err: ", err)
					return
				}
			}
		case <-waitChan:
			//fmt.Printf("normal quit job ppid:%d\n", cmd.Process.Pid)
		}
	}()

	if err := cmd.Wait(); err != nil {
		//fmt.Printf("timeout kill job ppid:%s\n%s\n", b.String(), err.Error())
		em := err.Error()
		// 超时退出，返回调用失败
		if strings.Contains(em, "signal: killed") {
			return "", errors.Wrap(err, "exec timeout")
		}
		// 未超时，被执行程序主动退出
		if i.errPrint {
			fmt.Printf("stdout: %s\n", stdout.String())
			fmt.Printf("stderr: %s\n", stderr.String())
		}
		if len(stderr.String()) != 0 {
			return stdout.String(), errors.New(stderr.String())
		}
		return stdout.String() + "\n" + stderr.String(), errors.Wrap(err, "")
	}

	out = stdout.String()
	// 唤起正常退出
	waitChan <- struct{}{}

	return
}

func addEnv(cmd *exec.Cmd, envs map[string]string) {
	if len(cmd.Env) == 0 {
		cmd.Env = []string{}
	}
	for key, value := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
}

// Exe 根据运行环境中是否具有sh或bash来执行命令command，如果执行后有stderr或执行90s超时，返回err；否则返回执行结果到string
func Exe(dir string, name string, arg ...string) (string, error) {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(90)*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	cmd = exec.Command(name, arg...)
	if dir != "" {
		cmd.Dir = dir
	}

	var ob bytes.Buffer
	var eb bytes.Buffer
	// 使子进程拥有自己的 pgid，等同于子进程的 pid
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = &ob
	cmd.Stderr = &eb

	err := cmd.Start()
	if err != nil {
		return "", errors.Wrapf(err, "start cmd")
	}

	waitChan := make(chan struct{}, 1)
	defer close(waitChan)

	// 超时杀掉进程组 或正常退出
	go func() {
		select {
		case <-ctx.Done():
			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
				return
			}
		case <-waitChan:
		}
	}()

	err = cmd.Wait()
	if err != nil {
		em := err.Error()
		// 超时退出，返回调用失败
		if strings.Contains(em, "signal: killed") {
			return "", errors.Wrapf(err, "exec [%s] timeout after 90s", cmd.String())
		}
		// 未超时，被执行程序主动退出
		if eb.String() == "" {
			return ob.String(), nil
		} else {
			return "", errors.Wrapf(err, "exec cmd with stderr [%s]", eb.String())
		}
	}

	obs := ob.String()
	ebs := eb.String()
	// 唤起正常退出
	waitChan <- struct{}{}
	if len(ebs) != 0 {
		return "", errors.New(ebs)
	}
	return obs, nil
}
