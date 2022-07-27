# 三、拦截器

基本市面上所有框架都是支持拦截器的，gRPC也不例外。grpc支持客户端和服务端拦截以下是Java和go语言的实现。

## 一、Java服务端拦截

### 1、创建一个proto文件进行演示

```protobuf
syntax = "proto3";//标识 proto版本 建议使用proto3
package orderservice;//proto包名 避免命名冲突，也可以作为引入其他proto文件时使用
option java_package = "com.example.orderservice" ;//生成的类将带有此包名，不指定则使用package
option cc_generic_services = true;
option go_package = "./pb";
option java_outer_classname = "OrderPbEntity";

message GetOrderListReq{
  string userId = 1;
}
message GetOrderListRes{
  repeated string orderIds = 1;
}

message SearchOrderListReq{
  string orderName = 1;
}
message SearchOrderListRes{
  repeated string orderIds = 1;
}

service OrderService {
  rpc getOrderList(GetOrderListReq) returns (GetOrderListRes);
}
```

### 2、实现服务端逻辑

```java
package com.example.orderservice.impl;

import com.example.orderservice.OrderPbEntity;
import com.example.orderservice.OrderServiceGrpc;
import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import io.grpc.stub.StreamObserver;

import java.util.HashMap;
import java.util.List;

public class OrderServiceImpl extends OrderServiceGrpc.OrderServiceImplBase {

    private HashMap<String, List<String>> data = Maps.newHashMap();

    @Override
    public void getOrderList(OrderPbEntity.GetOrderListReq request, StreamObserver<OrderPbEntity.GetOrderListRes> responseObserver) {
        if (data.isEmpty()){
            initData();
        }
        String userId = request.getUserId();
        List<String> strings = data.get(userId);
        responseObserver.onNext(OrderPbEntity.GetOrderListRes.newBuilder().addAllOrderIds(strings).build());
        responseObserver.onCompleted();
    }

    private void initData() {
        List<String> data1 = Lists.newArrayList();
        data1.add("order_1");
        data1.add("order_2");
        data1.add("order_3");
        data1.add("order_4");
        data.put("1",data1);
        data.put("2",data1);

    }
}
```

### 3、创建拦截器实现

- `OrderMsgServerCall`

```java
package com.example.orderservice.interceptors;

import io.grpc.ForwardingServerCall;
import io.grpc.MethodDescriptor;
import io.grpc.ServerCall;
import java.util.logging.Logger;

public class OrderMsgServerCall<ReqT, RespT> extends ForwardingServerCall.SimpleForwardingServerCall<ReqT, RespT>  {

    private static final Logger logger = Logger.getLogger(OrderMsgServerCall.class.getName());

    OrderMsgServerCall(ServerCall<ReqT, RespT> delegate) {
        super(delegate);
    }

    @Override
    protected ServerCall<ReqT, RespT> delegate() {
        return super.delegate();
    }

    @Override
    public MethodDescriptor<ReqT, RespT> getMethodDescriptor() {
        return super.getMethodDescriptor();
    }

    @Override
    public void sendMessage(RespT message) {
        logger.info("Message from Service -> Client : " + message);
        super.sendMessage(message);
    }
}
```

- `OrderMsgServerCallListener`

```java
package com.example.orderservice.interceptors;

import com.example.orderservice.OrderPbEntity;
import io.grpc.ForwardingServerCallListener;
import io.grpc.ServerCall;

import java.util.logging.Logger;

public class OrderMsgServerCallListener<R> extends ForwardingServerCallListener<R> {
    private static final Logger logger = Logger.getLogger(OrderMsgServerCallListener.class.getName());

    private final ServerCall.Listener<R> delegate;

    OrderMsgServerCallListener(ServerCall.Listener<R> delegate) {
        this.delegate = delegate;
    }

    @Override
    protected ServerCall.Listener<R> delegate() {
        return delegate;
    }

    @Override
    public void onMessage(R message) {
        if(message instanceof OrderPbEntity.GetOrderListReq){
            OrderPbEntity.GetOrderListReq req = (OrderPbEntity.GetOrderListReq)message;
            String userId = req.getUserId();
            logger.info("userId:"+userId);
        }
        logger.info("Message Received from Client -> Service " + message);
        super.onMessage(message);
    }
}
```

- `ParamsInterceptor`

```java
package com.example.orderservice.interceptors;

import io.grpc.*;

import java.util.logging.Logger;

public class ParamsInterceptor implements ServerInterceptor {
    Logger logger = Logger.getLogger(this.getClass().getName());
    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(ServerCall<ReqT, RespT> call, Metadata headers, ServerCallHandler<ReqT, RespT> next) {
        logger.info("======= [Server Interceptor] : Remote Method Invoked - " + call.getMethodDescriptor().getFullMethodName());
        ServerCall<ReqT, RespT> serverCall = new OrderMsgServerCall<>(call);
        return new OrderMsgServerCallListener<>(next.startCall(serverCall, headers));
    }
}
```

### 4、服务启动注入拦截器

```java
package com.example.orderservice;

import com.example.orderservice.impl.OrderServiceImpl;
import com.example.orderservice.interceptors.ParamsInterceptor;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.protobuf.services.ProtoReflectionService;

import java.io.IOException;

public class App {
    public String getGreeting() {
        return "Hello World!";
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        System.out.println(new App().getGreeting());

        Server server = ServerBuilder.forPort(10082).intercept(new ParamsInterceptor()).addService(new OrderServiceImpl())
                .addService(ProtoReflectionService.newInstance())
                .build();
        server.start();
        server.awaitTermination();
    }
}
```

## 二、Java客户端拦截

- 客户端启动类 直接注入拦截器

```java
/*
 * This Java source file was generated by the Gradle 'init' task.
 */
package com.examples.orderservice;

import com.example.orderservice.OrderPbEntity;
import com.example.orderservice.OrderServiceGrpc;
import io.grpc.*;

public class App {

    public static void main(String[] args) {
        ManagedChannel managedChannel = ManagedChannelBuilder.forAddress("localhost",10082)
                .usePlaintext()
                .build();
        Channel channel = ClientInterceptors.intercept(managedChannel, new ClientInterceptor() {
            @Override
            public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(MethodDescriptor<ReqT, RespT> method, CallOptions callOptions, Channel next) {
                System.out.println("method:"+method.getFullMethodName());
                return next.newCall(method,callOptions);
            }
        });
        OrderServiceGrpc.OrderServiceBlockingStub orderServiceBlockingStub = OrderServiceGrpc.newBlockingStub(channel);

        OrderPbEntity.GetOrderListReq req = OrderPbEntity.GetOrderListReq.newBuilder().setUserId("1")
                .build();
        OrderPbEntity.GetOrderListRes orderList = orderServiceBlockingStub.getOrderList(req);
        for (String orderId : orderList.getOrderIdsList()) {
            System.out.println("orderId:"+orderId);
        }

    }
}
```

## 三、Go服务端拦截器

grpc-go除了支持客户端和服务端拦截外还区分开了流式拦截和普通一元拦截，由于处理方式基本一致这里只展示普通rpc请求的拦截

```go
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
```

## 四、客户端拦截

```go
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
```