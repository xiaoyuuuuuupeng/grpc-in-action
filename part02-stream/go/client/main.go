package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-in-action/part02-stream/go/client/pb"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:10081", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserInfoServiceClient(conn)

	GetUserInfoRes(c)
	ReportUserBehavior(c)   //客户端stream
	SupplyUserChangeInfo(c) //服务端stream
	ExchangeUserInfo(c)     //客户端服务端stream

}

func ExchangeUserInfo(c pb.UserInfoServiceClient) {
	info, err := c.ExchangeUserInfo(context.Background())
	if err != nil {
		panic(err)
	}
	for i := 0; i < 5; i++ {
		req := &pb.ExchangeUserInfoReq{
			Id:  "1",
			Lng: 12.12,
			Lat: 32.32,
		}
		err = info.Send(req)
		if err != nil {
			panic(err)
		}
	}
	lisChanClose := make(chan bool)

	go listenExchangeUserInfo(info, lisChanClose)

	if err := info.CloseSend(); err != nil {
		log.Fatal(err)
	}

	lisChanClose <- true

}

func listenExchangeUserInfo(info pb.UserInfoService_ExchangeUserInfoClient, listenClose chan bool) {
	for {
		data, err := info.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		log.Printf("exchange Data : %v", data)
	}
	log.Printf("bool:%v", <-listenClose)
}

func SupplyUserChangeInfo(c pb.UserInfoServiceClient) {
	req := &pb.SupplyUserChangeInfoReq{
		Id: "1",
	}
	res, err := c.SupplyUserChangeInfo(context.Background(), req)
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}
	for {
		data, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		log.Printf("data:%v", data)
	}
}

func ReportUserBehavior(c pb.UserInfoServiceClient) {
	req := &pb.ReportUserBehaviorReq{
		Id:       "1",
		Behavior: "起床",
	}

	req1 := &pb.ReportUserBehaviorReq{
		Id:       "1",
		Behavior: "洗漱",
	}

	req2 := &pb.ReportUserBehaviorReq{
		Id:       "1",
		Behavior: "上班",
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	bhClient, err := c.ReportUserBehavior(context)
	if err != nil {
		panic(err)
	}

	if err = bhClient.Send(req); err != nil {
		panic(err)
	}

	if err = bhClient.Send(req); err != nil {
		panic(err)
	}

	if err = bhClient.Send(req1); err != nil {
		panic(err)
	}

	if err = bhClient.Send(req2); err != nil {
		panic(err)
	}
	bhRes, err := bhClient.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", bhClient, err, nil)
	}
	log.Printf("Update Orders Res : %s", bhRes)

}

func GetUserInfoRes(c pb.UserInfoServiceClient) {

	req := &pb.GetUserInfoReq{Id: "1"}
	res, err := c.GetUserInfo(context.Background(), req)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("GetUserInfoRes:%v", res)
}
