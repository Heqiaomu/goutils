package eventbus

import (
	log "github.com/Heqiaomu/glog"
	"sync"
)

// Bus
// Subscribe 订阅消息，特定event 使用特定Handler完成
// Publish 发布特定event 执行Handler
// UnSubscribe 退订消息
type Bus struct {
	event    chan Event
	special  chan *OnceEvent
	handlers map[string]func(Event)
	lock     sync.Mutex
}
type OnceEvent struct {
	event  Event
	handle func(Event)
}

func NewBus() *Bus {
	return &Bus{
		event:    make(chan Event, 1024),
		special:  make(chan *OnceEvent, 1024),
		handlers: make(map[string]func(Event), 0),
		lock:     sync.Mutex{},
	}
}
func (b *Bus) Subscribe(topic string, handle func(Event)) {
	b.lock.Lock()
	defer b.lock.Unlock()
	// 直接覆盖
	b.handlers[topic] = handle
}
func (b *Bus) HandleEvent() {
	for true {
		select {
		case e := <-b.event:
			h, ok := b.handlers[e.GetTopic()]
			if ok {
				go h(e)
			} else {
				log.Errorf("Fail\u001B[33mTopic[%v] no registered handler\u001B[0m.", e.GetTopic())
			}
		case once := <-b.special:
			log.Infof("Invoke\u001B[33m Special topic(%s)\u001B[0m.", once.event.GetTopic())
			go once.handle(once.event)
		default:
			// 没有发布时间不需要执行
		}
	}
}
func (b *Bus) Publish(e Event) {
	b.event <- e
	log.Debugf("\u001B[31mMessage's Topic(%s) published\u001B[0m.", e.GetTopic())
}
func (b *Bus) OncePublish(e Event, handle func(Event)) {
	b.special <- &OnceEvent{
		event:  e,
		handle: handle,
	}
	log.Debugf("\u001B[31mSpecial topic(%s) published\u001B[0m.", e.GetTopic())
}

type Event interface {
	GetTopic() string
	// Register(bus *Bus)
}

func Register(bus *Bus, topic string, handle func(event Event)) {
	bus.Subscribe(topic, handle)
}
