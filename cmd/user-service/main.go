// cmd/user-service/main.go
package main

import (
	"context"
	"fmt"
	"helloworld-go/internal/bootstrap"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	ginEngine, grpcServer, err := bootstrap.Build()
	if err != nil {
		log.Fatal(err)
	}

	// 启动 HTTP
	go func() {
		if err := ginEngine.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	// 启动 gRPC
	lis, _ := net.Listen("tcp", ":50051")
	grpcSrv := grpc.NewServer()
	grpcSrv.RegisterService(&grpcServer.ServiceDesc, grpcServer)
	fmt.Println("gRPC server running on :50051")
	grpcSrv.Serve(lis)
}
