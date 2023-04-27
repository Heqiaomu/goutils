package websocket

import (
	"testing"

	log "github.com/Heqiaomu/glog"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	log.Logger()
}

func TestServer(t *testing.T) {
	Convey("websocket client_Start", t, func() {
		Convey("websocket client test ", func() {
			//c,err := NewWsClient("ws://127.0.0.1:8080/ws?agentId=123")
			//if err != nil{
			//	return
			//}
			//time.Sleep(100*time.Second)
			//fmt.Println(c)
			//go func() {
			//	for{
			//		j:=simplejson.New()
			//		j.Set(CMD,"abc")
			//		j.Set(DATA,"1111")
			//
			//		b,_:= j.MarshalJSON()
			//		c.SendMessage(string(b))
			//		//rsp ,err := c.PostMessage(string(b))
			//		//if err != nil{
			//		//	fmt.Println(err)
			//		//	continue
			//		//}
			//		//fmt.Println("收到",rsp)
			//		time.Sleep(3*time.Second)
			//	}
			//
			//}()
			//time.Sleep(1000*time.Second)
			//c.Stop()
			ShouldBeNil(nil)
		})
	})
}

func TestNewWsClient(t *testing.T) {
	type args struct {
		coreAddr string
		opts     []WsClientOption
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Client
		wantErr bool
	}{
		{
			"TestNewWsClient_1",
			args{
				"127.0.0.1",
				[]WsClientOption{SetIsAlwaysRetry(false)},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWsClient(tt.args.coreAddr, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWsClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// gotC.onMessage("Xxx")
		})
	}
}

func TestNewWsClientEx(t *testing.T) {
	type args struct {
		coreAddr string
		onMsg    func(string)
		opts     []WsClientOption
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Client
		wantErr bool
	}{
		{
			"TestNewWsClient_1",
			args{
				"127.0.0.1",
				func(string) {},
				[]WsClientOption{SetIsAlwaysRetry(false)},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWsClientEx(tt.args.coreAddr, tt.args.onMsg, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWsClientEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
