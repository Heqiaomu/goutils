package server

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	Convey("TestServer_Start", t, func() {
		Convey("no server", func() {
			PushMsgToServer("new server", "1")
		})

		Convey("Should be start", func() {
			s, err := NewServerEx("new server", 1,
				func(msg string, num int) error {
					return nil
				},
				1*time.Second,
				[]ScheduledTask{func(num int) {
					fmt.Println(3, num)
				}, func(num int) {
					fmt.Println(4, num)
				}})

			s1, err := NewServer("new server", 3, 2, func(msg string, num int) error {
				fmt.Println(num, msg)
				return nil
			}, func(num int) {
				fmt.Println("循环中")
			})

			s.StartEx()

			PushMsgToServer("new server", time.Now().String())

			time.Sleep(1 * time.Second)
			//s.Stop()
			time.Sleep(1 * time.Millisecond)
			DestoryServer("new server")

			ShouldBeError(err)
			ShouldBeNil(s1)
		})

		Convey("go server", func() {
			s, _ := NewSvr("test-go-server",
				func(msg string, num int) error {
					fmt.Println(msg, num)
					return nil
				},
				[]TimedTask{
					{
						Task: func(num int) {
							fmt.Println(time.Now(), num)
						},
						Time: 15 * time.Second,
					},
					{
						Task: func(num int) {
							fmt.Println(time.Now(), num)
						},
						Time: 5 * time.Second,
					},
				})
			s.Go()
			go func() {
				PushMsgToServer("test-go-server", "abc")
			}()
			time.Sleep(1 * time.Second)

		})
	})
}
