package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-in-action/part03-interceptor/go/client/pb"
	"log"
)
func orderUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Pre-processor phase
	log.Println("Method : " + method)
	// Invoking the remote method
	err := invoker(ctx, method, req, reply, cc, opts...)
	// Post-processor phase
	log.Println(reply)
	return err
}



func main() {

	dial, err := grpc.Dial("localhost:10082", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(orderUnaryClientInterceptor),
		)
	if err != nil{panic(err)}

	client := pb.NewOrderServiceClient(dial)
	list, err := client.GetOrderList(context.Background(), &pb.GetOrderListReq{UserId: "go_client_userId"})
	if err != nil {
		log.Fatalf("err:%v",err)
	}
	log.Printf("list:%v",list.OrderIds)
}
