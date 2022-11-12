package main

import (
	"context"
	"fmt"
	"go-grpc/pb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTestServer
}

func (*server) Unary(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	fmt.Println("キャッチした: ", req.GetMessage())
	res := &pb.TestResponse{}
	res.Message = "起きた"
	time.Sleep(2 * time.Second)

	fmt.Println("送る値: ", res.GetMessage())
	return res, nil
}

func (*server) ServerStreaming(req *pb.TestRequest, stream pb.Test_ServerStreamingServer) error {
	res := &pb.TestResponse{}
	fmt.Println("キャッチした: ", req.GetMessage())
	for i := 0; i < 5; i++ {
		res.Message = "起きろ"
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		fmt.Println("送る値: ", res.GetMessage())
		time.Sleep(1 * time.Second)
	}

	return nil

}

func (*server) ClientStreaming(stream pb.Test_ClientStreamingServer) error {

	var result string = ""
	for {
		req, err := stream.Recv()
		name := req.GetName()
		fmt.Println("キャッチした: ", req.GetName())
		result += "Hello " + name + ", "
		if err == io.EOF {
			return stream.SendAndClose(&pb.TestResponse{
				Message: result,
			})
		}
		if err != nil {
			return err
		}
	}
}

func (*server) BidirectionalStreaming(stream pb.Test_BidirectionalStreamingServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println("送られた値: ", req.GetMessage())
		if req.GetMessage() == "あなたたちは?" {
			talks := []string{"とても多いですね。", "私たちはGo1", "Go2", "Go3", "Go4", "Go5です"}

			for _, talk := range talks {
				fmt.Println("送る値:", talk)
				res := &pb.TestResponse{
					Message: talk,
				}
				time.Sleep(1 * time.Second)
				err = stream.Send(res)
				if err != nil {
					return err
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:6565")
	if err != nil {
		log.Fatalf("失敗 %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTestServer(s, &server{})

	fmt.Println("server is running...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("失敗 %v", err)
	}
}
