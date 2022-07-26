package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-in-action/part02-stream/go/server/impl"
	"grpc-in-action/part02-stream/go/server/pb"
	"log"
	"net"
)

func main() {

	conn, err := net.Listen("tcp", "localhost:10081")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	pb.RegisterUserInfoServiceServer(server, &impl.StreamUserInfoServiceImpl{})
	log.Println("服务器启动，监听在：10081")
	// Register reflection service on gRPC server.
	reflection.Register(server)
	err = server.Serve(conn)
	if err != nil {
		panic(err)
	}
}
