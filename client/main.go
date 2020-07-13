package main

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-plugins/transport/rabbitmq/v2"

	br "github.com/micro/go-plugins/broker/rabbitmq/v2"
	proto "github.com/zhsso/greeter-demo/proto"

	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/transport"

	"github.com/micro/go-micro/v2/registry/etcd"
)

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf("[wrapper] client request to service: %s method: %s\n", req.Service(), req.Method())
	return c.Client.Call(ctx, req, rsp)
}

// 返回一个包装过的客户端
func LogClientWrap(c client.Client) client.Client {
	return &clientWrapper{c}
}

func main() {
	ts := rabbitmq.NewTransport(func(options *transport.Options) {
		options.Addrs = []string{"amqp://192.168.50.116:5672"}
	})
	brk := br.NewBroker(func(options *broker.Options) {
		options.Addrs = []string{"amqp://192.168.50.116:5672"}
	})
	reg := etcd.NewRegistry(registry.Addrs("192.168.50.116:2379"))
	// Create a new service
	service := micro.NewService(
		micro.Client(client.NewClient()),
		micro.Server(server.NewServer()),
		micro.Broker(brk),
		micro.Registry(reg),
		micro.Transport(ts),
		micro.Name("greeter.client"),
		micro.WrapClient(LogClientWrap),
	)
	// Initialise the client and parse command line flags
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())

	_, err := greeter.Hello(context.TODO(), &proto.Request{Name: "John Supper"})
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 20; i++ {
		//fmt.Println("begin hello", i)
		// Call the greeter
		rsp, err := greeter.Hello(context.TODO(), &proto.Request{Name: fmt.Sprintf("John %d", i)})
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Minute)
		} else {

			// Print response
			fmt.Println(rsp.Greeting)
			//time.Sleep(time.Second * 1)
		}
	}
}
