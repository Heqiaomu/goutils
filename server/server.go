package server

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const QueueLen = 1000

var Servers map[string]*Server = make(map[string]*Server)

var ServersNew = sync.Map{}

type GoState int

const (
	// server is stopped
	Stopped GoState = iota
	// server is started
	Started
)

type MsgAction func(msg string, num int) error
type ScheduledTask func(num int)

type Server struct {
	Name                      string
	Queue                     chan string
	State                     GoState
	ActionGoroutineNum        int
	MsgActionerExit           []chan int
	MsgActioner               MsgAction
	ScheduledTaskGoroutineNum int
	ScheduledTaskerExit       []chan int
	ScheduledTaskers          []ScheduledTask
	ScheduleTime              time.Duration
	TimedTasks                []TimedTask
}

// Deprecated: Use NewSvr instead.
func NewServer(serverName string, actionGoroutineNum int, scheduledTaskGoroutineNum int, action MsgAction, schedule ScheduledTask) (server *Server, err error) {
	_, ok := ServersNew.Load(serverName)
	if ok {
		return nil, fmt.Errorf("the same server serverName already exists:%s", serverName)
	}
	//s := Servers[serverName]
	//if s != nil {
	//	return nil, fmt.Errorf("the same server serverName already exists:%s", serverName)
	//}
	server = &Server{
		Name:                      serverName,
		ActionGoroutineNum:        actionGoroutineNum,
		Queue:                     make(chan string, QueueLen),
		MsgActionerExit:           []chan int{},
		MsgActioner:               action,
		ScheduledTaskGoroutineNum: scheduledTaskGoroutineNum,
		ScheduledTaskerExit:       []chan int{},
	}
	var mschedules []ScheduledTask
	for i := 0; i < scheduledTaskGoroutineNum; i++ {
		mschedules = append(mschedules, schedule)
	}
	server.ScheduledTaskers = mschedules

	for i := 0; i < server.ActionGoroutineNum; i++ {
		server.MsgActionerExit = append(server.MsgActionerExit, make(chan int))
	}
	for i := 0; i < server.ScheduledTaskGoroutineNum; i++ {
		server.ScheduledTaskerExit = append(server.ScheduledTaskerExit, make(chan int))
	}

	server.State = Stopped
	ServersNew.Store(serverName, server)
	//Servers[serverName] = server
	return server, nil
}

// Deprecated: Use NewSvr instead.
// 简单版本
func NewServerEx(serverName string, actionGoroutineNum int, action MsgAction, scheduleTime time.Duration, schedules []ScheduledTask) (server *Server, err error) {
	//s := Servers[serverName]
	//if s != nil {
	//	return nil, fmt.Errorf("the same server serverName already exists:%s", serverName)
	//}
	_, ok := ServersNew.Load(serverName)
	if ok {
		return nil, fmt.Errorf("the same server serverName already exists:%s", serverName)
	}
	server = &Server{
		Name:                      serverName,
		ActionGoroutineNum:        actionGoroutineNum,
		Queue:                     make(chan string, QueueLen),
		MsgActionerExit:           []chan int{},
		MsgActioner:               action,
		ScheduledTaskGoroutineNum: len(schedules),
		ScheduledTaskerExit:       []chan int{},
		ScheduledTaskers:          schedules,
		ScheduleTime:              scheduleTime,
	}
	for i := 0; i < server.ActionGoroutineNum; i++ {
		server.MsgActionerExit = append(server.MsgActionerExit, make(chan int))
	}
	for i := 0; i < server.ScheduledTaskGoroutineNum; i++ {
		server.ScheduledTaskerExit = append(server.ScheduledTaskerExit, make(chan int))
	}

	server.State = Stopped
	//Servers[serverName] = server
	ServersNew.Store(serverName, server)
	return server, nil
}

type TimedTask struct {
	Task ScheduledTask
	Time time.Duration
}

func NewSvr(serverName string, action MsgAction, timedTasks []TimedTask) (server *Server, err error) {
	//s := Servers[serverName]
	//if s != nil {
	//	return nil, fmt.Errorf("the same server serverName already exists:%s", serverName)
	//}
	_, ok := ServersNew.Load(serverName)
	if ok {
		return nil, fmt.Errorf("the same server serverName already exists:%s", serverName)
	}
	server = &Server{
		Name:                      serverName,
		ActionGoroutineNum:        runtime.NumCPU() * 2,
		Queue:                     make(chan string, QueueLen),
		MsgActionerExit:           []chan int{},
		MsgActioner:               action,
		ScheduledTaskGoroutineNum: len(timedTasks),
		ScheduledTaskerExit:       []chan int{},
		TimedTasks:                timedTasks,
	}
	// 若没有消息接受需要处理，则将线程数量=0
	if action == nil {
		server.ActionGoroutineNum = 0
	}
	for i := 0; i < server.ActionGoroutineNum; i++ {
		server.MsgActionerExit = append(server.MsgActionerExit, make(chan int))
	}
	for i := 0; i < server.ScheduledTaskGoroutineNum; i++ {
		server.ScheduledTaskerExit = append(server.ScheduledTaskerExit, make(chan int))
	}

	server.State = Stopped
	//Servers[serverName] = server
	ServersNew.Store(serverName, server)
	return server, nil
}

func DestoryServer(serverName string) {
	//s := Servers[serverName]
	load, ok := ServersNew.Load(serverName)
	if !ok {
		return
	}
	s := load.(*Server)
	s.Stop()
	close(s.Queue)
	ServersNew.Delete(serverName)
	//delete(Servers, serverName)
}

func PushMsgToServer(serverName string, msg string) error {
	//s := Servers[serverName]
	load, ok := ServersNew.Load(serverName)
	if !ok {
		return fmt.Errorf("server name doesn't exist:%s", serverName)
	}
	s := load.(*Server)
	s.ReceiveMsg(msg)
	return nil
}

// Deprecated: Use Go instead.
func (server *Server) StartEx() {
	server.Start(server.ScheduleTime)
}

// Deprecated: Use Go instead.
// 开启服务，定时任务时间间隔：ScheduleTime
func (server *Server) Start(scheduleTime time.Duration) {

	for i := 0; i < server.ActionGoroutineNum; i++ {
		number := i
		go func() {
			//fmt.Printf("[%s][%d] MsgActioner goroutine 启动\n", server.Name, number)
			for {
				select {
				case msg := <-server.Queue:
					{
						if server.MsgActioner != nil {
							err := server.MsgActioner(msg, number)
							if err != nil {
								fmt.Println(err)
							}
						}
					}
				case <-server.MsgActionerExit[number]:
					//fmt.Printf("[%s][%d] MsgActioner goroutine 退出\n", server.Name, number)
					return
				}
			}
		}()
	}
	for j := 0; j < server.ScheduledTaskGoroutineNum; j++ {
		number := j
		go func() {
			//fmt.Printf("[%s][%d] ScheduledTasker goroutine 启动\n", server.Name, number)
			server.ScheduledTaskers[number](number)
			for {
				select {
				case <-server.ScheduledTaskerExit[number]:
					//fmt.Printf("[%s][%d] ScheduledTasker goroutine 退出\n", server.Name, number)
					return
				case <-time.After(scheduleTime):
					server.ScheduledTaskers[number](number)
				}
			}
		}()
	}

	server.State = Started

	return
}

// 开启服务，定时任务时间间隔：ScheduleTime
func (server *Server) Go() {

	for i := 0; i < server.ActionGoroutineNum; i++ {
		number := i
		go func() {
			//fmt.Printf("[%s][%d] MsgActioner goroutine 启动\n", server.Name, number)
			for {
				select {
				case msg := <-server.Queue:
					{
						if server.MsgActioner != nil {
							err := server.MsgActioner(msg, number)
							if err != nil {
								fmt.Println(err)
							}
						}
					}
				case <-server.MsgActionerExit[number]:
					//fmt.Printf("[%s][%d] MsgActioner goroutine 退出\n", server.Name, number)
					return
				}
			}
		}()
	}
	for j := 0; j < server.ScheduledTaskGoroutineNum; j++ {
		number := j
		go func() {
			//fmt.Printf("[%s][%d] ScheduledTasker goroutine 启动\n", server.Name, number)
			server.TimedTasks[number].Task(number)
			for {
				select {
				case <-server.ScheduledTaskerExit[number]:
					//fmt.Printf("[%s][%d] ScheduledTasker goroutine 退出\n", server.Name, number)
					return
				case <-time.After(server.TimedTasks[number].Time):
					server.TimedTasks[number].Task(number)
				}
			}
		}()
	}

	server.State = Started

	return
}

func (server *Server) Stop() {
	if server.State == Started {
		for i := 0; i < server.ActionGoroutineNum; i++ {
			server.MsgActionerExit[i] <- 1
		}
		for i := 0; i < server.ScheduledTaskGoroutineNum; i++ {
			server.ScheduledTaskerExit[i] <- 1
		}
		server.State = Stopped
	}
}

func (server *Server) ReceiveMsg(msg string) {
	server.Queue <- msg
}
