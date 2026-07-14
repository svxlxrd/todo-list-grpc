package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/svxlxrd/todo-list-grpc/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type todoServer struct {
	pb.UnimplementedTodoServiceServer
	mu    sync.RWMutex
	tasks []*pb.Task
}

func (s *todoServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &pb.Task{
		Id:    fmt.Sprintf("%d", len(s.tasks)+1),
		Title: req.Title,
		Done:  false,
	}
	s.tasks = append(s.tasks, task)

	log.Printf("Task created: %s", task.Title)
	return task, nil
}

func (s *todoServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		if task.Id == req.Id {
			return task, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "task not found")
}

func (s *todoServer) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*pb.Task
	for _, t := range s.tasks {
		if req.OnlyPending && t.Done {
			continue
		}
		result = append(result, t)
	}
	return &pb.ListTasksResponse{Tasks: result}, nil
}

func (s *todoServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.tasks {
		if t.Id == req.Id {
			t.Done = req.Done
			return t, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "task not found")
}

func (s *todoServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks {
		if t.Id == req.Id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return &pb.DeleteTaskResponse{Success: true}, nil
		}
	}
	return &pb.DeleteTaskResponse{Success: false}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &todoServer{})

	log.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
