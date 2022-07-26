package main

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"grpc-in-action/part04-metedata/go/client/pb"
	"log"
)

func main() {

	conn,err:=grpc.Dial("localhost:10083",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	req := &pb.PayOrderReq{OrderId: "order_2"}

	client := pb.NewPayServiceClient(conn)
	token := uuid.New().String()

	md := metadata.Pairs(
		"token",token)
	context := metadata.NewOutgoingContext(context.Background(),md)
	res,err:=client.PayOrder(context,req)
	if err != nil {
		panic(err)
	}
	log.Printf("retCode:%d,retMsg:%s",res.RetCode,res.RetMsg)
}
