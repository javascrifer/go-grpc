package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	c := greetpb.NewGreetServiceClient(cc)
	callGreet(c)
	callGreetManyTime(c)
	callLongGreet(c)
}

func callGreet(c greetpb.GreetServiceClient) {
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

func callGreetManyTime(c greetpb.GreetServiceClient) {
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

func callLongGreet(c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}

	for i := 0; i < 10; i++ {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nikolajus",
				LastName:  "Lebedenko",
			},
		}
		stream.Send(req)
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while reading response of LongGreet RPC: %v", err)
	}
	log.Printf("response from LongGreet RPC: %v", res.GetResult())
}
