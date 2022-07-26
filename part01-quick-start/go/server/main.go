package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-in-action/part1/go/server/pb"
	"grpc-in-action/part1/go/server/serviceimpl"
	"log"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", "localhost:10080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ser := &serviceimpl.UserInfoServiceImpl{}
	pb.RegisterUserInfoServiceServer(s, ser)
	log.Printf("start")
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
