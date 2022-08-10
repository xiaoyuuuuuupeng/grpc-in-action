package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-in-action/part05-errorhanle/go/server/impl"
	"grpc-in-action/part05-errorhanle/go/server/pb"
	"log"
	"net"
)

func main() {

	listen, err := net.Listen("tcp", "localhost:10085")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	pb.RegisterBlendServiceServer(server,&impl.BlendService{})
	reflection.Register(server)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
