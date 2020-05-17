package main

import (
	"context"
	"fmt"
	"time"

	proto "github.com/zhsso/greeter-demo/proto"

	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/transport/nats/v2"
)

func main() {
	reg := consul.NewRegistry()
	ts := nats.NewTransport()

	// Create a new service
	service := micro.NewService(
		micro.Client(client.NewClient()),
		micro.Transport(ts),
		micro.Registry(reg),
		micro.Name("greeter.client"),
	)
	// Initialise the client and parse command line flags
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())

	for i := 0; i < 100; i++ {
		fmt.Println("begin hello", i)
		// Call the greeter
		rsp, err := greeter.Hello(context.TODO(), &proto.Request{Name: "John"})
		if err != nil {
			fmt.Println(err)
		}

		// Print response
		fmt.Println(rsp.Greeting)
		time.Sleep(time.Second * 1)
	}
}
