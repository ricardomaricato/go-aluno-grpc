package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ricardomaricato/go-aluno-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	AddUsers(client)
	// AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}
	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}
	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "r1",
			Name:  "Ricardo 1",
			Email: "ricardo@ricardo.com",
		},
		&pb.User{
			Id:    "r2",
			Name:  "Ricardo 2",
			Email: "ricardo@ricardo.com",
		},
		&pb.User{
			Id:    "r3",
			Name:  "Ricardo 3",
			Email: "ricardo@ricardo.com",
		},
		&pb.User{
			Id:    "r4",
			Name:  "Ricardo 4",
			Email: "ricardo@ricardo.com",
		},
		&pb.User{
			Id:    "r2",
			Name:  "Ricardo 5",
			Email: "ricardo@ricardo.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}
