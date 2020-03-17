package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/javascrifer/go-grpc/internal/pkg/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

// Greet is implementing Unary gRPC API call
func (s *server) Greet(
	ctx context.Context,
	req *greetpb.GreetRequest,
) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	if firstName == "" || lastName == "" {
		err := status.Errorf(
			codes.InvalidArgument,
			"received empty first name or the last name",
		)
		return nil, err
	}

	// Uncomment in order to fake some computing and deadline handling
	// time.Sleep(50 * time.Millisecond)
	// if ctx.Err() == context.DeadlineExceeded {
	// log.Println("client cancelled Greet RPC call")
	// err := status.Errorf(codes.Canceled, "client cancelled the request")
	// return nil, err
	// }

	res := &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hello %s %s!", firstName, lastName),
	}

	return res, nil
}

// GreetManyTimes is implementing Server streaming gRPC API call
func (s *server) GreetManyTimes(
	req *greetpb.GreetManyTimesRequest,
	stream greetpb.GreetService_GreetManyTimesServer,
) error {
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

// LongGreet is implementing Client streaming gRPC API call
func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked")
	fullNames := []string{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &greetpb.LongGreetResponse{
				Result: fmt.Sprintf("Hello %s!", strings.Join(fullNames, ", ")),
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatalf("error while reading stream in LongGreet RPC: %v\n", err)
			return err
		}

		fmt.Printf("LongGreet received request %v\n", req)
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		fullNames = append(fullNames, fmt.Sprintf("%s %s", firstName, lastName))
	}
}

// GreetEveryone is implementing BiDi streaming gRPC API call
func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading stream in GreetEveryone RPC: %v\n", err)
			return err
		}

		fmt.Printf("GreetEveryone received request %v\n", req)
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		res := &greetpb.GreetEveryoneResponse{
			Result: fmt.Sprintf("Hello %s %s!", firstName, lastName),
		}

		if err := stream.Send(res); err != nil {
			log.Fatalf("error while writing to stream in GreetEveryone RPC: %v\n", err)
			return err
		}
	}
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
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
