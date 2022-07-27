package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-in-action/part03-interceptor/go/server/pb"
	"log"
	"net"
)

type OrderServiceImpl struct {
	pb.UnimplementedOrderServiceServer
}

func (impl *OrderServiceImpl) GetOrderList(ctx context.Context, req *pb.GetOrderListReq) (*pb.GetOrderListRes, error) {

	userId := req.UserId
	log.Printf("userId:%s",userId)
	orderIds:=[]string{"order_go_1","order_go_2","order_go_3"}
	return &pb.GetOrderListRes{OrderIds: orderIds},nil
}
// Server :: Unary Interceptor
func orderUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Pre-processing logic
	// Gets info about the current RPC call by examining the args passed in
	log.Println("======= [Server Interceptor] ", info.FullMethod)
	log.Printf(" Pre Proc Message : %s", req)
	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)
	// Post processing logic
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}


func main() {

	listen, err := net.Listen("tcp", "localhost:10082")
	if err != nil{
		panic(err)
	}
	server :=grpc.NewServer(grpc.ChainUnaryInterceptor(orderUnaryServerInterceptor))
	reflection.Register(server)
	pb.RegisterOrderServiceServer(server,&OrderServiceImpl{})
	if err = server.Serve(listen) ;err != nil {
		panic(err)
	}
}
