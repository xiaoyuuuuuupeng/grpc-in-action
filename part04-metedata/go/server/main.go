package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"grpc-in-action/part04-metedata/go/server/pb"
	"log"
	"net"
	"time"
)

type PayServiceImpl struct {
	pb.UnimplementedPayServiceServer
}

func (pay *PayServiceImpl)PayOrder(context context.Context, req *pb.PayOrderReq) (*pb.PayOrderRes, error) {

	md,boo := metadata.FromIncomingContext(context)
	if boo{
		log.Printf("token:%s",md["token"])
	}
	log.Printf("req:orderId:%s",req.OrderId)

	// Creating and sending a header.
	header := metadata.New(map[string]string{"location": "San Jose", "timestamp": time.Now().Format(time.StampNano)})
	grpc.SendHeader(context, header)
	return &pb.PayOrderRes{
		RetCode: 0,
		RetMsg:  "支付成功",
	},nil
}

func main() {

	listen, err := net.Listen("tcp", "localhost:10083")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	pb.RegisterPayServiceServer(server,&PayServiceImpl{})
	reflection.Register(server)
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}


}
