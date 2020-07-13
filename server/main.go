package main

import (
	"context"
	"fmt"
	"time"

	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/transport"
	br "github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/micro/go-plugins/transport/rabbitmq/v2"
	proto "github.com/zhsso/greeter-demo/proto"

	"github.com/micro/go-micro/v2/registry/etcd"
)

// handler 包装
func HandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		fmt.Printf("[HandlerWrapper] [%v] server request: %s : %s\n", time.Now(), req.Service(), req.Method())
		return fn(ctx, req, rsp)
	}
}

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Greeting = "Hello " + req.Name
	println("return hello!")
	return nil
}

func main() {
	ts := rabbitmq.NewTransport(func(options *transport.Options) {
		options.Addrs = []string{"amqp://192.168.50.116:5672"}
	})
	brk := br.NewBroker(func(options *broker.Options) {
		options.Addrs = []string{"amqp://192.168.50.116:5672"}
	})
	reg := etcd.NewRegistry(registry.Addrs("192.168.50.116:2379"))
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Server(server.NewServer()),
		micro.Client(client.NewClient()),
		micro.Broker(brk),
		micro.Registry(reg),
		micro.Transport(ts),
		micro.Name("greeter"),
		micro.WrapHandler(HandlerWrapper),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	if err := proto.RegisterGreeterHandler(service.Server(), new(Greeter)); err != nil {
		panic(err)
	}

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
