package wserver

import (
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	Convey("websocket server", t, func() {
		Convey("websocket server test ", func() {

			//viper.Set("websocket.port",9999)
			//NewServer(func(conn *websocket.Conn, agentId string) {
			//
			//}, func(conn *websocket.Conn, agentId string, message []byte) {
			//	fmt.Println(message)
			//}, func(conn *websocket.Conn, agentId string) {
			//
			//})
			//Start()
			//time.Sleep(1000*time.Second)
			ShouldBeNil(nil)
		})
	})
}

func Test_ws(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws(tt.args.w, tt.args.r)
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		onConnect OnConnnectFunc
		onMessage OnMessageFunc
		onClose   OnCloseFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"TestNewServer",
			args{
				func(conn *websocket.Conn, agentId string) {

				},
				func(conn *websocket.Conn, agentId string, message []byte) {

				},
				func(conn *websocket.Conn, agentId string) {

				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewServer(tt.args.onConnect, tt.args.onMessage, tt.args.onClose)
			NewServerWithMux(http.NewServeMux(), tt.args.onConnect, tt.args.onMessage, tt.args.onClose)
		})
	}
}
