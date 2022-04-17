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
	fmt.Println("HelloReply was invoked")
	res := &pb.TestResponse{}
	res.Message = req.GetName()

	return res, nil
}

func (*server) ServerStreaming(req *pb.TestRequest, stream pb.Test_ServerStreamingServer) error {
	fmt.Println("ServerStreaming was invoked")
	res := &pb.TestResponse{}

	for i := 0; i < 5; i++ {

		res.Message = req.GetName()
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}

	return nil

}

func (*server) ClientStreaming(stream pb.Test_ClientStreamingServer) error {
	fmt.Println("Upload was invoked")

	var sum int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			message := fmt.Sprintf("DONE: sum = %d", sum)
			return stream.SendAndClose(&pb.TestResponse{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		fmt.Println(req.GetName())
		sum = 15
	}
}

func (*server) BidirectionalStreaming(stream pb.Test_BidirectionalStreamingServer) error {
	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		data := req.GetName()
		log.Printf("received data: %v", data)
		size += len(data)

		res := &pb.TestResponse{
			Message: fmt.Sprintf("received %v bytes", size),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}
func (*server) AStreaming(req *pb.ARequest, stream pb.Test_AStreamingServer) error {
	fmt.Println("ServerStreaming was invoked")
	res := &pb.AResponse{}

	for i := 0; i < 5; i++ {
		res.Message = req.GetName()
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}

	return nil

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
