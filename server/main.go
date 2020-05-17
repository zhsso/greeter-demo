package main

import (
	"context"
	"fmt"

	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/transport/nats/v2"

	proto "github.com/zhsso/greeter-demo/proto"

	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	reg := consul.NewRegistry()
	ts := nats.NewTransport()
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Server(server.NewServer()),
		micro.Transport(ts),
		micro.Registry(reg),
		micro.Name("greeter"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
