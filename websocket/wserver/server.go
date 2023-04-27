package wserver

import (
	"fmt"
	log "github.com/Heqiaomu/glog"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"net/http"
	"sync"
)

type OnConnnectFunc func(conn *websocket.Conn, agentId string)
type OnMessageFunc func(conn *websocket.Conn, agentId string, message []byte)
type OnCloseFunc func(conn *websocket.Conn, agentId string)

var OnConnect OnConnnectFunc
var OnMessage OnMessageFunc
var OnClose OnCloseFunc

func SetOnConnect(function OnConnnectFunc) {
	OnConnect = function
}

// 设置消息接受处理回调方法
func SetOnMessage(function OnMessageFunc) {
	OnMessage = function
}

// 设置链接断开处理回调方法
func SetOnClose(function OnCloseFunc) {
	OnClose = function
}

var Conns = &sync.Map{}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ConnInfo struct {
	Conn *websocket.Conn
	Lock *sync.Mutex
}

func (c *ConnInfo) SendMsg(msg string) (err error) {
	c.Lock.Lock()
	err = c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
	c.Lock.Unlock()
	return err
}

func ws(w http.ResponseWriter, r *http.Request) {

	v := r.URL.Query()
	agentId := v.Get("agentId")
	if agentId == "" {
		log.Warnf("Currently the agentID is null")
		return
	}

	if err := SaveWsConn(agentId); err != nil {
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Currently upgrade:", err)
		return
	}

	ci := &ConnInfo{
		Conn: c,
		Lock: &sync.Mutex{},
	}
	Conns.Store(agentId, ci)
	OnConnect(c, agentId)

	defer c.Close()

	c.SetPingHandler(func(appData string) (err error) {

		ci.Lock.Lock()
		err = c.WriteMessage(websocket.PongMessage, []byte(appData))
		ci.Lock.Unlock()

		if err != nil {
			log.Errorf("Currently 发送pong失败", err)
			return err
		}

		return nil
	})

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			ciptr, ok := Conns.Load(agentId)
			if ok == true {
				if ciptr.(*ConnInfo).Conn == c {
					log.Warnf("Currently agentId=%s 断开, err=%v", agentId, err)
					Conns.Delete(agentId)
				}
			}
			OnClose(c, agentId)
			break
		}
		OnMessage(c, agentId, message)
	}
}

func SaveWsConn(agentId string) error {
	if agentId == "" {
		return fmt.Errorf("the agentID is null")
	}

	_, ok := Conns.Load(agentId)
	if ok == true {
		log.Warnf("Currently the same AgentID[%s] already exists, close the conn", agentId)
		//ciptr.(*ConnInfo).Conn.Close()

		////等待连接断开
		//for i := 0;i<50;i++{
		//	_,ok := Conns.Load(agentId)
		//	if ok == false{
		//		fmt.Println("相同的agentId已经断开")
		//		break
		//	}
		//	time.Sleep(100*time.Millisecond)
		//}
	}

	return nil
}

var Server *http.Server = nil

func NewServer(onConnect OnConnnectFunc, onMessage OnMessageFunc, onClose OnCloseFunc) {
	SetOnConnect(onConnect)
	SetOnMessage(onMessage)
	SetOnClose(onClose)

	http.HandleFunc("/ws", ws)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("websocket"))))
}

func NewServerWithMux(httpmux *http.ServeMux, onConnect OnConnnectFunc, onMessage OnMessageFunc, onClose OnCloseFunc) {
	SetOnConnect(onConnect)
	SetOnMessage(onMessage)
	SetOnClose(onClose)

	httpmux.HandleFunc("/ws", ws)
	httpmux.Handle("/wstest/", http.StripPrefix("/wstest/", http.FileServer(http.Dir("./websocket/"))))
}

func Start() {
	go func() {
		Server = &http.Server{Addr: "0.0.0.0:" + viper.GetString("websocket.port")}
		Server.ListenAndServe()
	}()
}

func Stop() {
	Server.Close()
}

func SendMsg(agentId string, msg string) (err error) {
	ci, ok := Conns.Load(agentId)
	if ok == false {
		return fmt.Errorf("连接已经不存在")
	}
	err = ci.(*ConnInfo).SendMsg(msg)
	return
}
