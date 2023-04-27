package websocket

import (
	"context"
	"fmt"
	log "github.com/Heqiaomu/glog"
	"github.com/Heqiaomu/goutil/server"
	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"
	"time"
)

const CMD string = "Cmd"
const UUID string = "Uuid"
const DATA string = "Data"

type Client struct {
	Conn       *Connect       // 当前ws连接句柄
	WsServer   *server.Server // 服务句柄
	PongSignal chan string
	ReqSignal  map[string]chan string
	OnMsg      func(string)
	wso        WsClientConf
}

type WsClientConf struct {
	IsAlwaysRetry bool
}

type WsClientOption func(c *WsClientConf)

func SetIsAlwaysRetry(v bool) WsClientOption {
	return func(c *WsClientConf) {
		c.IsAlwaysRetry = v
	}
}

func NewWsClient(coreAddr string, opts ...WsClientOption) (c *Client, err error) {

	c = &Client{PongSignal: make(chan string), ReqSignal: map[string]chan string{}, OnMsg: nil, wso: WsClientConf{}}

	// 设置客户端配置项
	for _, opt := range opts {
		opt(&c.wso)
	}

	c.Conn, err = NewConn(coreAddr, c.onMessage, c.onPong, c.wso)
	if err != nil {
		return nil, err
	}

	c.WsServer, err = server.NewServer(fmt.Sprintf("ws-client-%s", uuid.New().String()), 1, 1,
		// 发送消息
		func(msg string, num int) error {
			j, err := simplejson.NewJson([]byte(msg))
			if err != nil {
				return err
			}
			cmd, err := j.Get(CMD).String()
			if err != nil {
				return err
			}
			data, err := j.Get(DATA).MarshalJSON()
			if err != nil {
				return err
			}

			if cmd == "ping" {
				if err = c.Conn.Ping(string(data)); err != nil {
					log.Errorf("Currently ws ping 失败: %v", err)
				}
			} else if cmd == "text" {
				if err = c.Conn.SendMsg(string(data)); err != nil {
					log.Errorf("Currently ws send text 失败:", err)
				}
			}
			return nil
		},
		// 维护心跳
		func(num int) {
			j := simplejson.New()
			j.Set(CMD, "ping")
			j.Set(DATA, uuid.New().String())

			// 发送ping包
			jStr, _ := j.MarshalJSON()
			c.WsServer.ReceiveMsg(string(jStr))
			select {
			case <-time.After(3 * time.Second): // 3s钟未收到pong消息，认为链接已经断开
				{
					log.Warn("Currently ws连接中断，进行重连")
					c.Conn.StopClient()
					// 重新创建conn
					c.Conn, err = NewConn(coreAddr, c.onMessage, c.onPong, c.wso)
					if err != nil {
						log.Errorf("Currently new conn err: %v", err)
					}
				}
			case <-c.PongSignal:
				{
				}
			}
			return
		},
	)

	if err != nil {
		return nil, err
	}
	c.WsServer.Start(3 * time.Second)

	return
}

func NewWsClientEx(coreAddr string, onMsg func(string), opts ...WsClientOption) (c *Client, err error) {

	c = &Client{PongSignal: make(chan string), ReqSignal: map[string]chan string{}, OnMsg: onMsg, wso: WsClientConf{}}

	for _, opt := range opts {
		opt(&c.wso)
	}

	c.Conn, err = NewConn(coreAddr, c.onMessage, c.onPong, c.wso)
	if err != nil {
		return nil, err
	}

	c.WsServer, err = server.NewServer(fmt.Sprintf("ws-client-%s", uuid.New().String()), 1, 1,
		// 发送消息
		func(msg string, num int) error {
			j, err := simplejson.NewJson([]byte(msg))
			if err != nil {
				return err
			}
			cmd, err := j.Get(CMD).String()
			if err != nil {
				return err
			}
			data, err := j.Get(DATA).MarshalJSON()
			if err != nil {
				return err
			}

			if cmd == "ping" {
				c.Conn.Ping(string(data))
			} else if cmd == "text" {
				c.Conn.SendMsg(string(data))
			}
			return nil
		},
		// 维护心跳
		func(num int) {
			j := simplejson.New()
			j.Set(CMD, "ping")
			j.Set(DATA, uuid.New().String())

			// 发送ping包
			jStr, _ := j.MarshalJSON()
			c.WsServer.ReceiveMsg(string(jStr))
			select {
			case <-time.After(3 * time.Second): // 3s钟未收到pong消息，认为链接已经断开
				{
					log.Warn("Currently ws连接中断，进行重连")
					c.Conn.StopClient()
					// 重新创建conn
					c.Conn, err = NewConn(coreAddr, c.onMessage, c.onPong, c.wso)
					if err != nil {
						log.Errorf("Currently new conn err: %v", err)
					}
				}
			case <-c.PongSignal:
				{
				}
			}
			return
		},
	)

	if err != nil {
		return nil, err
	}
	c.WsServer.Start(3 * time.Second)

	return
}

func (c *Client) onMessage(message string) {
	j, err := simplejson.NewJson([]byte(message))
	if err != nil {
		if c.OnMsg != nil {
			c.OnMsg(message)
		}
		return
	}
	uuid, err := j.Get(UUID).String()
	t, ok := c.ReqSignal[uuid]
	if ok == true {
		t <- message
	} else {
		if c.OnMsg != nil {
			c.OnMsg(message)
		}
	}
	return
}

func (c *Client) onPong(message string) {
	c.PongSignal <- message
	return
}

// SendMessage 用于发送不接受回包的消息
func (c *Client) SendMessage(msg string) (err error) {
	// 构造uuid消息
	uuid := uuid.New().String()
	reqJson, err := simplejson.NewJson([]byte(msg))
	reqJson.Set(UUID, uuid)

	j := simplejson.New()
	j.Set(CMD, "text")
	j.Set(DATA, reqJson)

	mb, err := j.MarshalJSON()
	if err != nil {
		return
	}
	c.WsServer.ReceiveMsg(string(mb))
	return
}

// PostMessage 用于发送要立马收到回包的消息，返回服务端返回的消息
func (c *Client) PostMessage(req string) (rsp string, err error) {

	// 构造uuid消息
	uuid := uuid.New().String()
	reqJson, err := simplejson.NewJson([]byte(req))
	reqJson.Set(UUID, uuid)

	// 构造req消息chan
	t := make(chan string)
	c.ReqSignal[uuid] = t

	// 构造发送消息体，并发送
	j := simplejson.New()
	j.Set(CMD, "text")
	j.Set(DATA, reqJson)
	mb, err := j.MarshalJSON()
	if err != nil {
		return
	}
	c.WsServer.ReceiveMsg(string(mb))

	select {
	case rsp = <-t:
		{
		}
	case <-time.After(10 * time.Second):
		{
			err = fmt.Errorf("消息发送超时,message:%s", req)
		}
	}
	// 删除req消息chan
	delete(c.ReqSignal, uuid)

	return
}

// ConnMessage 用于发送要持续收到回包的消息，返回服务端返回的消息
func (c *Client) ConnMessage(ctxWithCancel context.Context, req string, f func(string) error) (err error) {

	// 构造uuid消息
	uuid := uuid.New().String()
	reqJson, err := simplejson.NewJson([]byte(req))
	reqJson.Set(UUID, uuid)

	// 构造req消息chan
	t := make(chan string, 1000)
	c.ReqSignal[uuid] = t

	// 构造发送消息体，并发送
	j := simplejson.New()
	j.Set(CMD, "text")
	j.Set(DATA, reqJson)
	mb, err := j.MarshalJSON()
	if err != nil {
		return
	}
	c.WsServer.ReceiveMsg(string(mb))

	for {
		select {
		case rsp, ok := <-t:
			if !ok {
				goto end
			}
			err := f(rsp)
			if err != nil {
				log.Error(err.Error())
			}
		case <-ctxWithCancel.Done():
			goto end
		}
	}
end:
	//删除req消息chan
	delete(c.ReqSignal, uuid)

	return nil
}

func (c *Client) Stop() {
	server.DestoryServer(c.WsServer.Name)
	//c.WsServer.Stop()
	c.Conn.StopClient()
}
