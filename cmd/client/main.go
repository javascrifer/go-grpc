package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/javascrifer/go-grpc/internal/pkg/calculatorpb"
	"github.com/javascrifer/go-grpc/internal/pkg/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Initializing client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer cc.Close()

	callGreet(cc)
	callGreetManyTime(cc)
	callAdd(cc)
}

func callGreet(cc *grpc.ClientConn) {
	c := greetpb.NewGreetServiceClient(cc)
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nikolajus",
			LastName:  "Lebedenko",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("response from Greet RPC: %v", res.Result)
}

func callGreetManyTime(cc *grpc.ClientConn) {
	c := greetpb.NewGreetServiceClient(cc)
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nikolajus",
			LastName:  "Lebedenko",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream in GreetManyTimes RPC: %v", err)
		}
		log.Printf("response from GreetManyTimes RPC: %v", message.GetResult())
	}
}

func callAdd(cc *grpc.ClientConn) {
	c := calculatorpb.NewCalculatorServiceClient(cc)
	req := &calculatorpb.CalculatorRequest{
		X: 1,
		Y: 2,
	}
	res, err := c.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Add RPC: %v", err)
	}
	log.Printf("response from Add RPC: %v", res.Sum)
}
