package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-in-action/part1/go/client/pb"
	"log"
)

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:10080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserInfoServiceClient(conn)
	req := &pb.GetUserInfoReq{Id: "1"}
	res, err := c.GetUserInfo(context.Background(), req)
	if err != nil {
		log.Fatalf("err:%+v", err)
	} else {
		log.Printf("res:%+v", res)
	}

}
