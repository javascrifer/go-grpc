package main

import (
	"context"
	"fmt"
	"log"

	"github.com/javascrifer/go-grpc/internal/pkg/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Initializing server")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	callGreet(c)
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
