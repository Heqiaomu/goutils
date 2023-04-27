package server

import (
	"github.com/google/uuid"
	"sync"
)

// Factory 工厂，工人的集合
type Factory struct {
	Workers sync.Map
}

// NewFactory 建造一个工厂
func NewFactory() (factory *Factory) {
	return &Factory{}
}

// Destroy 摧毁一个工厂，工人全部停工
func (factory *Factory) Destroy() {
	factory.Workers.Range(func(key, value interface{}) bool {
		value.(*Worker).Stop()
		return true
	})
}

// GetWorker 根据ID获取工人
func (factory *Factory) GetWorker(ID string) *Worker {
	v, ok := factory.Workers.Load(ID)
	if ok {
		return v.(*Worker)
	}
	return nil
}

// NewWorker 招聘一个新工人
func (factory *Factory) NewWorker(labels []string) (worker *Worker) {
	worker = NewWorker(labels)
	factory.Workers.Store(worker.ID, worker)
	return
}

// FireWorker 解雇一个工人同时停止工人的工作
func (factory *Factory) FireWorker(ID string) {
	w := factory.GetWorker(ID)
	if w != nil {
		w.Stop()
		factory.Workers.Delete(ID)
	}
}

// Worker 工人
type Worker struct {
	ID     string    // 工人编号
	Doing  bool      // 工人干活状态，true：正在干活 false：停止工作
	Exit   chan bool // 停止工作信号
	Labels []string
}

// NewWorker 生成一个工人
func NewWorker(labels []string) (worker *Worker) {
	return &Worker{ID: uuid.New().String(), Exit: make(chan bool), Doing: false, Labels: labels}
}

// WorkerAction 工人的工作 task：工作内容 exit：停止工作信号
type WorkerAction func(task interface{}, exit chan bool)

// Start 工人开始工作
func (worker *Worker) Start(wa WorkerAction, task interface{}) {
	worker.Doing = true
	go func() {
		wa(task, worker.Exit)
	}()
}

// Stop 工人停止工作
func (worker *Worker) Stop() {
	worker.Exit <- true
	worker.Doing = false
}
