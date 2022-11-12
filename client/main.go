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
	Unary(client)
	// ServerStreaming(client)
	// ClientStreaming(client)
	// BidirectionalStreaming(client)
}

func Unary(client pb.TestClient) {

	req := &pb.TestRequest{}
	req.Message = "起きろ"
	fmt.Println("送る値:", req.GetMessage())
	res, err := client.Unary(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res)
}

func ServerStreaming(client pb.TestClient) error {
	req := &pb.TestRequest{}
	req.Message = "明日起こして"
	fmt.Println("送る値: ", req.GetMessage())
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
		log.Println("返ってきた値: ", reply.Message)
	}
	req2 := &pb.TestRequest{}
	req2.Message = "起きた"
	fmt.Println("送る値: ", req2.GetMessage())
	client.Unary(context.Background(), req2)
	return nil
}

func ClientStreaming(client pb.TestClient) error {
	stream, err := client.ClientStreaming(context.Background())
	if err != nil {
		return err
	}
	values := []string{"太郎", "二郎", "三郎", "四郎", "五郎"}
	for _, value := range values {
		fmt.Println("送る値:", value)
		if err := stream.Send(&pb.TestRequest{
			Name: value,
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
			talks := []string{"私たちは", "太郎", "二郎", "三郎", "四郎", "五郎", "あなたたちは?"}
			for _, talk := range talks {
				fmt.Println("送る値:", talk)
				req := &pb.TestRequest{
					Message: talk,
				}

				time.Sleep(1 * time.Second)
				err = stream.Send(req)
			}
			time.Sleep(1 * time.Second)
			return
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
			fmt.Println("返ってきたよ: ", res.GetMessage())

			if err != nil {
				log.Fatalln(err)
				break
			}
		}
		close(ch)
	}()
	<-ch
}
