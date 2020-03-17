package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/javascrifer/go-grpc/internal/pkg/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
	callGreetEveryone(c)
}

// Unary
func callGreet(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nikolajus",
			LastName:  "Lebedenko", // change to "" in order request to fail
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Fatalf("[%v] error while calling Greet RPC: %v\n", s.Code(), s.Message())
		} else {
			log.Fatalf("big error while calling Greet RPC: %v\n", err)
		}
		return
	}
	log.Printf("response from Greet RPC: %v\n", res.Result)
}

// Server streaming
func callGreetManyTime(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nikolajus",
			LastName:  "Lebedenko",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v\n", err)
		return
	}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream in GreetManyTimes RPC: %v\n", err)
			return
		}
		log.Printf("response from GreetManyTimes RPC: %v\n", message.GetResult())
	}
}

// Client streaming
func callLongGreet(c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v\n", err)
		return
	}

	for i := 0; i < 10; i++ {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nikolajus",
				LastName:  "Lebedenko",
			},
		}
		fmt.Printf("sending message for LongGreet RPC stream: %v\n", req)
		stream.Send(req)
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while reading response of LongGreet RPC: %v\n", err)
		return
	}
	log.Printf("response from LongGreet RPC: %v\n", res.GetResult())
}

// BiDi streaming
func callGreetEveryone(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone RPC: %v\n", err)
		return
	}

	waitc := make(chan bool)

	// Send stream
	go func() {
		for i := 0; i < 10; i++ {
			req := &greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Nikolajus",
					LastName:  "Lebedenko",
				},
			}
			fmt.Printf("sending message for GreetEveryone RPC stream: %v\n", req)
			stream.Send(req)
			time.Sleep(time.Duration(rand.Intn(int(time.Second) * 5)))
		}
		stream.CloseSend()
	}()

	// Receive stream
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream in GreetEveryone RPC: %v\n", err)
				break
			}
			fmt.Printf("response from GreetEveryone RPC: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
	fmt.Println("GreetEveryone finished")
}
