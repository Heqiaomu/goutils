package eventbus

import (
	log "github.com/Heqiaomu/glog"
	"strconv"
	"testing"
	"time"
)

type ChainCreateTest struct {
	someMessage string
	result      bool
}

func (c *ChainCreateTest) GetTopic() string {
	return "create-chain-test"
}
func PrintEventMessage(event Event) {
	create, ok := event.(*ChainCreateTest)
	if ok {
		log.Infof("\u001B[31mEvent message: %s\u001B[0m.", create.someMessage)
	} else {
		log.Errorf("\u001B[%v] is not *ChainCreateTest\u001B[0m.", event)
	}
	time.Sleep(5 * time.Second)
}
func SpecialEventHandle(event Event) {
	create, ok := event.(*ChainCreateTest)
	if ok {
		log.Infof("\u001B[31mSpecial Handle event: %s\u001B[0m.", create.someMessage)
	} else {
		log.Errorf("\u001B[%v] is not *ChainCreateTest\u001B[0m.", event)
	}
	time.Sleep(5 * time.Second)
}

func TestEventBus(t *testing.T) {
	log.Logger()
	bus := NewBus()
	// 开始处理event
	go bus.HandleEvent()
	create0 := &ChainCreateTest{someMessage: "主线程完成"}
	// 注册普通的事件
	Register(bus, create0.GetTopic(), PrintEventMessage)
	// 发布
	bus.Publish(create0)

	for i := 0; i < 3; i++ {
		inner := strconv.Itoa(i)
		event := &ChainCreateTest{
			someMessage: "线程" + inner + "完成",
		}
		go func() {
			// some logic
			bus.Publish(event)
		}()
	}
	// 执行特殊的事件处理
	bus.OncePublish(create0, SpecialEventHandle)
	time.Sleep(30 * time.Second)
}
