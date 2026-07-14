package main

import (
	"context"
	"log"
	"time"

	pb "github.com/svxlxrd/todo-list-grpc/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// вызовы
	res, err := c.CreateTask(ctx, &pb.CreateTaskRequest{
		Title: "Купить молоко",
	})
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	log.Printf("Response from server: Task ID: %s, Title: %s", res.Id, res.Title)
}
