package http

import (
	"github.com/fztcjjl/tiger/pkg/middleware/gin/trace"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler() http.Handler {
	route := gin.New()

	route.Use(trace.Trace())
	route.GET("/hello", sayHello)
	return route
}

func sayHello(ctx *gin.Context) {
	//cli := client.NewClient(
	//	"srv.hello",
	//	client.Registry(etcd.NewRegistry(registry.Addrs("127.0.0.1:2379"))),
	//	client.GrpcDialOption(grpc.WithInsecure()),
	//	client.Interceptors(
	//		grpc_opentracing.UnaryClientInterceptor(),
	//	),
	//)
	//
	//grpcClient := pb.NewGreeterClient(cli.GetConn())
	//req := pb.HelloRequest{Name: "John"}
	//rsp, err := grpcClient.SayHello(trace.ContextWithSpan(ctx), &req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Greeting: %s", rsp.Message)
	//ctx.Writer.WriteString(rsp.Message)
}
