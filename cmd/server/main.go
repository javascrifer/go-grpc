package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/javascrifer/go-grpc/internal/pkg/calculatorpb"
	"github.com/javascrifer/go-grpc/internal/pkg/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	res := &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hello %s %s!", firstName, lastName),
	}

	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %s %s %v!", firstName, lastName, i),
		}
		stream.Send(res)
		time.Sleep(time.Second)
	}

	return nil
}

func (s *server) Add(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResposne, error) {
	fmt.Printf("Add function was invoked with %v\n", req)

	x := req.GetX()
	y := req.GetY()
	res := &calculatorpb.CalculatorResposne{
		Sum: x + y,
	}

	return res, nil
}

func main() {
	fmt.Println("Initializing server")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := &server{}
	greetpb.RegisterGreetServiceServer(s, srv)
	calculatorpb.RegisterCalculatorServiceServer(s, srv)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
