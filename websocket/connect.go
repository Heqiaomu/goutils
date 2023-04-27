package websocket

import (
	"fmt"
	log "github.com/Heqiaomu/glog"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Connect struct {
	Conn *websocket.Conn
	Lock *sync.Mutex
}

type OnMessageFunc func(message string)
type OnPongFunc func(message string)

func NewConn(coreAddr string, messageFunc OnMessageFunc, pongFunc OnPongFunc, conf WsClientConf) (conn *Connect, err error) {
	var c *websocket.Conn
	for i := 0; i < 3; i++ {
		// 创建conn句柄
		c, _, err = websocket.DefaultDialer.Dial(coreAddr, nil)
		if err != nil {
			log.Errorf("Currently websocket dail to %s, error : %v", coreAddr, err)
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		break
	}
	if err != nil && !conf.IsAlwaysRetry {
		return nil, err
	}

	// 构造Connet对象
	conn = &Connect{Conn: c, Lock: &sync.Mutex{}}

	log.Debugf("ws 连接成功:%s", coreAddr)

	// 设置pong回调
	c.SetPongHandler(func(appData string) error {
		pongFunc(appData)
		return nil
	})

	// 开启接受消息的goroutine
	go func() {
		for {
			_, p, err := c.ReadMessage()
			if err != nil {
				log.Debugf("ws client is closed:", err.Error())
				break
			}
			// 设置收到消息的回调
			messageFunc(string(p))
		}
	}()
	return
}

// StopClient 停止ws client
func (c *Connect) StopClient() {
	if c == nil || c.Conn == nil {
		log.Info("Currently conn is nil don't need to stop")
		return
	}
	c.Lock.Lock()
	if c.Conn != nil {
		c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.Conn.Close()
	}
	c.Lock.Unlock()
}

// SendMsg 发送msg
func (c *Connect) SendMsg(msg string) error {
	if c == nil || c.Conn == nil {
		return fmt.Errorf("conn is closed,can not send msg")
	}
	var err error
	c.Lock.Lock()
	err = c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
	c.Lock.Unlock()
	return err
}

// 发送ping包
func (c *Connect) Ping(msg string) error {
	if c == nil || c.Conn == nil {
		return fmt.Errorf("conn is closed,can not ping")
	}
	var err error
	c.Lock.Lock()
	err = c.Conn.WriteMessage(websocket.PingMessage, []byte(msg))
	c.Lock.Unlock()
	return err
}
