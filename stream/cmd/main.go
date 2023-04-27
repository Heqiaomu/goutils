package main

import (
	"fmt"
	"github.com/Heqiaomu/goutil/stream/stream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
	"time"
)

type service struct {
	stream.UnimplementedStreamServer
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	stream.RegisterStreamServer(s, &service{})
	s.Serve(listen)
}

// GetStream 服务端 单向流
func (s *service) GetStream(req *stream.StreamRequest, res stream.Stream_GetStreamServer) error {
	i := 0
	for {
		i++
		err := res.Send(&stream.StreamResponse{Output: fmt.Sprintf("%v", time.Now().Unix())})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
		if i > 1000000 {
			break
		}
	}

	return nil
}

// PutStream 客户端 单向流
func (s *service) PutStream(cliStr stream.Stream_PutStreamServer) error {

	for {
		if tem, err := cliStr.Recv(); err == nil {
			log.Println(tem)
		} else {
			log.Println("break, err :", err)
			break
		}
	}

	cliStr.SendAndClose(&stream.StreamResponse{Output: "abcd"})

	return nil
}

//客户端服务端 双向流
func (s *service) AllStream(allStr stream.Stream_AllStreamServer) error {

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for i := 0; i < 5; i++ {
			allStr.Send(&stream.StreamResponse{Output: "server"})
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	for {
		data, err := allStr.Recv()
		if err != nil {
			fmt.Println("客户端不发数据了")
			break
		}
		log.Println(data)
	}

	wg.Wait()

	return nil
}
