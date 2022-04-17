package main

import (
	"context"
	"fmt"
	"go-grpc/pb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:6565", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("失敗: %v", err)
	}
	defer conn.Close()

	client := pb.NewTestClient(conn)
	// Unary(client)
	// ServerStreaming(client)
	// ClientStreaming(client)
	BidirectionalStreaming(client)
}

func Unary(client pb.TestClient) {
	req := &pb.TestRequest{}
	req.Name = "Taro"
	res, err := client.Unary(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res)
}

// func AStreaming(client pb.TestClient) error {
// 	req := &pb.ARequest{}
// 	req.Name = "Taro"
// 	stream, err := client.AStreaming(context.Background(), req)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	for {
// 		reply, err := stream.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		log.Println("これ：", reply.GetMessage())
// 	}
// 	return nil
// }

func ServerStreaming(client pb.TestClient) error {
	req := &pb.TestRequest{}
	req.Name = "Taro"
	stream, err := client.ServerStreaming(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Println("これ：", reply.GetMessage())
	}
	return nil
}

func ClientStreaming(client pb.TestClient) error {
	stream, err := client.ClientStreaming(context.Background())
	if err != nil {
		return err
	}
	values := []int32{1, 2, 3, 4, 5}
	for _, value := range values {
		fmt.Println("送る値:", value)
		if err := stream.Send(&pb.TestRequest{
			Name: "aa",
		}); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		time.Sleep(time.Second * 1)
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("結果: %v", reply)
	return nil
}

func BidirectionalStreaming(client pb.TestClient) {
	stream, err := client.BidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// request
	go func() {
		for {
			if err != nil {
				log.Fatalln(err)
			}
			req := &pb.TestRequest{Name: "Taro"}
			sendErr := stream.Send(req)
			if sendErr != nil {
				log.Fatalln(sendErr)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("received message: %v", res.GetMessage())
		}
		close(ch)
	}()
	<-ch
}
