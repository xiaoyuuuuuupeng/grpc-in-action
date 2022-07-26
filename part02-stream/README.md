# 二、流式传输

我们第一节学习的是简单一元数据传输，也就是在rpc中最常见的一个请求对应一个响应。但往往在实际开发中可能有这样的场景：用户app位置上报、服务端大文件下发等这种持续字节流的情况。Grpc贴心的针对这些情况设置了三种流式传输模式：**客户端流、服务端流、客户端服务端双向流传输**

## 1、客户端单向流

现在我们在proto文件中创建一个客户端上报流数据的rpc方法 - **reportUserBehavior**

```protobuf
syntax = "proto3";//标识 proto版本 建议使用proto3
package userinfoservice;//proto包名 避免命名冲突，也可以作为引入其他proto文件时使用
option java_package = "com.example.userinfoservice" ;//生成的类将带有此包名，不指定则使用package
option cc_generic_services = true;
option go_package = "./pb";
option java_outer_classname = "UserInfoEntity";

message GetUserInfoReq{
  string id = 1;
}
message GetUserInfoRes{
  string id = 1;
  string name = 2;
  int32 age = 3;
}

message ReportUserBehaviorReq{
  string id = 1;
  string behavior = 2;
}

message ReportUserBehaviorRes{
  int32 retCode = 1;
  string retMsg = 2;
  int64 receivedCount = 3;
}

service UserInfoService {
  rpc getUserInfo(GetUserInfoReq) returns (GetUserInfoRes);
  rpc reportUserBehavior(stream ReportUserBehaviorReq) returns (ReportUserBehaviorRes);
}
```

reportUserBehavior 模拟的是客户端上报用户行为的例子，在日常开发中，这种场景很常见，比如打点环节里可能需要将用户的一系列的行为操作上报。我们使用go和Java分别来实现这个案例：

- **go语言实现**

a、通过protoc生成代码并实现客户端代码

```go
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
	conn,err:= grpc.Dial("localhost:10081",grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserInfoServiceClient(conn)

	ReportUserBehavior(c)//客户端stream
}

func ReportUserBehavior(c pb.UserInfoServiceClient) {
	req:=&pb.ReportUserBehaviorReq{
		Id:       "1",
		Behavior: "起床",
	}

	req1:=&pb.ReportUserBehaviorReq{
		Id:       "1",
		Behavior: "洗漱",
	}

	req2:=&pb.ReportUserBehaviorReq{
		Id:       "1",
		Behavior: "上班",
	}

	context,cancel := context.WithTimeout(context.Background(),time.Second *10)

	defer cancel()
	bhClient ,err := c.ReportUserBehavior(context)
	if err != nil {
		panic(err)
	}

	if err= bhClient.Send(req);err!=nil{
		panic(err)
	}

	if err= bhClient.Send(req);err!=nil{
		panic(err)
	}

	if err= bhClient.Send(req1);err!=nil{
		panic(err)
	}

	if err= bhClient.Send(req2);err!=nil{
		panic(err)
	}
	bhRes, err := bhClient.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", bhClient, err, nil)
	}
	log.Printf("Update Orders Res : %s", bhRes)

}

```

b、go服务端支持

```go
package impl

import (
	"context"
	"errors"
	"fmt"
	"grpc-in-action/part02-stream/go/server/pb"
	"io"
	"log"
)

type StreamUserInfoServiceImpl struct {
	UserLocationTmp map[string][]UserLocTmp
	UserInfoData map[string]*UserInfo
	pb.UnimplementedUserInfoServiceServer
}

type UserLocTmp struct {
	Loc string
	Temperature string
}

type UserInfo struct {
	Id string
	Name string
	Age int32
}

func (server *StreamUserInfoServiceImpl) GetUserInfo(ctx context.Context, in *pb.GetUserInfoReq) (*pb.GetUserInfoRes, error)  {

	log.Printf("go stream server收到请求")
	//初始化一个map
	if server.UserInfoData == nil{
		server.initUserInfoData()
	}
	data:= server.UserInfoData[in.Id]
	if data == nil{
		return nil,errors.New("该id不存在")
	}
	res := &pb.GetUserInfoRes{
		Id:   data.Id,
		Name: data.Name,
		Age:  data.Age,
	}

	return res,nil
}

func (server *StreamUserInfoServiceImpl) ReportUserBehavior(stream pb.UserInfoService_ReportUserBehaviorServer) error  {

	res:= pb.ReportUserBehaviorRes{
		RetCode: 0,
		RetMsg: "收到收到~",
		ReceivedCount: 0,
	}
	for {
		behaviorData, err := stream.Recv()
		if err == io.EOF {
			// 结束符 客户端停止发送
			return stream.SendAndClose(&res)
		}
		if err != nil {
			return err
		}
		//此处有并发问题 只是简单演示
		res.ReceivedCount+=1
		log.Printf("用户 ID : %s -行为： %s", behaviorData.Id, behaviorData.Behavior)
	}
}

func (server *StreamUserInfoServiceImpl) initUserInfoData (){
	server.UserInfoData = make(map[string]*UserInfo)

	server.UserInfoData["1"] = &UserInfo{
		Id:   "1",
		Name: "小美",
		Age:  18,
	}
	server.UserInfoData["2"] = &UserInfo{
		Id:   "2",
		Name: "小刚",
		Age:  28,
	}
	server.UserInfoData["3"] = &UserInfo{
		Id:   "3",
		Name: "小王",
		Age:  20,
	}
}

func (server *StreamUserInfoServiceImpl) initUserLocationTmp() {
	server.UserLocationTmp = make(map[string][]UserLocTmp)
	var data []UserLocTmp
	for i := 0; i < 10; i++ {
		t := UserLocTmp{
			Loc: "北京",
			Temperature :"26.5",
		}
		data = append(data,t)
	}
	server.UserLocationTmp["1"] = data
	server.UserLocationTmp["2"] = data
	server.UserLocationTmp["3"] = data
}
```

- **Java语言实现**

a、通过Gradle插件生成代码并实现客户端代码

```java
package com.example.userinfoservice;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.stub.StreamObserver;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

public class App {

    private static final Logger logger = Logger.getLogger(App.class.getName());

    public static void main(String[] args) throws InterruptedException {
        ManagedChannel channel = ManagedChannelBuilder.forAddress("localhost", 10081)
                .usePlaintext()
                .build();

        UserInfoServiceGrpc.UserInfoServiceBlockingStub stub =
                UserInfoServiceGrpc.newBlockingStub(channel);

        GetUserInfo(stub);
        //客户端流
        UserInfoServiceGrpc.UserInfoServiceStub asyncStub = UserInfoServiceGrpc.newStub(channel);
        ReportUserBehavior(asyncStub);

        channel.shutdown();
    }

    private static void ReportUserBehavior(UserInfoServiceGrpc.UserInfoServiceStub stub) {
        CountDownLatch finishLatch = new CountDownLatch(1);
        StreamObserver<UserInfoEntity.ReportUserBehaviorRes> responseObserver = new StreamObserver<UserInfoEntity.ReportUserBehaviorRes>() {
            @Override
            public void onNext(UserInfoEntity.ReportUserBehaviorRes value) {
                logger.info("Report User Behavior RetCode: " + value.getRetCode() + " RetMsg: " + value.getRetMsg());
            }
            @Override
            public void onError(Throwable t) {
            }
            @Override
            public void onCompleted() {
                logger.info("此次流处理完成");
                finishLatch.countDown();
            }
        };
        StreamObserver<UserInfoEntity.ReportUserBehaviorReq> reportUserBehaviorReqStreamObserver = stub.reportUserBehavior(responseObserver);
        reportUserBehaviorReqStreamObserver.onNext(UserInfoEntity.ReportUserBehaviorReq.newBuilder().setId("1")
                .setBehavior("下班").build());
        reportUserBehaviorReqStreamObserver.onNext(UserInfoEntity.ReportUserBehaviorReq.newBuilder().setId("1")
                .setBehavior("回家").build());
        reportUserBehaviorReqStreamObserver.onNext(UserInfoEntity.ReportUserBehaviorReq.newBuilder().setId("1")
                .setBehavior("做饭").build());

        if (finishLatch.getCount() == 0) {
            logger.warning("RPC completed or errored before we finished sending.");
            return;
        }
        reportUserBehaviorReqStreamObserver.onCompleted();
        try {
            if (!finishLatch.await(10, TimeUnit.SECONDS)) {
                logger.warning("FAILED : Process orders cannot finish within 10 seconds");
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

    }

    private static void GetUserInfo(UserInfoServiceGrpc.UserInfoServiceBlockingStub stub) {

        UserInfoEntity.GetUserInfoRes res = stub.getUserInfo(
                UserInfoEntity.GetUserInfoReq.newBuilder()
                        .setId("1")
                        .build());
        logger.info("用户id: " + res.getId() + "\t用户姓名:"+ res.getName()
                +"\t用户年龄:"+ res.getAge());
    }
}
```

b、服务端支持

```java
package com.example.userinfoservice.impl;

import com.example.userinfoservice.UserInfoEntity;
import com.example.userinfoservice.UserInfoServiceGrpc;
import com.example.userinfoservice.entity.UserInfoData;
import com.google.common.collect.Maps;
import io.grpc.stub.StreamObserver;
import org.apache.commons.lang3.StringUtils;

import java.util.HashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class UserInfoServiceImpl extends UserInfoServiceGrpc.UserInfoServiceImplBase {

    private HashMap<String, UserInfoData> userinfoHashMap = Maps.newHashMap();

    @Override
    public void getUserInfo(UserInfoEntity.GetUserInfoReq request, StreamObserver<UserInfoEntity.GetUserInfoRes> responseObserver) {
        System.out.println("java服务收到请求");
        if (StringUtils.isEmpty(request.getId())){
            responseObserver.onError(new Exception("id不能为空"));
            responseObserver.onCompleted();
            return;
        }

        if (userinfoHashMap.isEmpty()){
            initUserInfoHashMap();
        }

        UserInfoData userInfoData = userinfoHashMap.get(request.getId());
        if (userInfoData == null){
            responseObserver.onError(new Exception("id不能为空"));
            responseObserver.onCompleted();
            return;
        }
        responseObserver.onNext(UserInfoEntity.GetUserInfoRes.newBuilder()
                .setId(userInfoData.getId())
                .setName(userInfoData.getName())
                .setAge(userInfoData.getAge())
                .build());
        responseObserver.onCompleted();
    }

    //客户端流
    @Override
    public StreamObserver<UserInfoEntity.ReportUserBehaviorReq> reportUserBehavior(StreamObserver<UserInfoEntity.ReportUserBehaviorRes> responseObserver) {
        AtomicInteger count = new AtomicInteger(0);
        UserInfoEntity.ReportUserBehaviorRes.Builder res = UserInfoEntity.ReportUserBehaviorRes.newBuilder()
                .setRetCode(0)
                .setRetMsg("收到~~");

        StreamObserver<UserInfoEntity.ReportUserBehaviorReq> reportUserBehaviorReqStreamObserver = new StreamObserver<UserInfoEntity.ReportUserBehaviorReq>() {
            @Override
            public void onNext(UserInfoEntity.ReportUserBehaviorReq value) {
                System.out.println(value.getId() + "\t" + value.getBehavior());
                count.incrementAndGet();
            }

            @Override
            public void onError(Throwable t) {
            }

            @Override
            public void onCompleted() {
                res.setReceivedCount(count.get());
                responseObserver.onNext(res.build());

                responseObserver.onCompleted();
            }
        };
        return reportUserBehaviorReqStreamObserver;
    }

    private void initUserInfoHashMap() {
        userinfoHashMap.put("1",new UserInfoData("1","小美",24));
        userinfoHashMap.put("2",new UserInfoData("2","小强",28));
        userinfoHashMap.put("3",new UserInfoData("3","小刚",30));
    }
}
```

## 2、服务端单向流

在proto文件中创建一个服务端下发流数据的rpc方法 - **supplyUserChangeInfo**

```protobuf
syntax = "proto3";//标识 proto版本 建议使用proto3
package userinfoservice;//proto包名 避免命名冲突，也可以作为引入其他proto文件时使用
option java_package = "com.example.userinfoservice" ;//生成的类将带有此包名，不指定则使用package
option cc_generic_services = true;
option go_package = "./pb";
option java_outer_classname = "UserInfoEntity";

message GetUserInfoReq{
  string id = 1;
}
message GetUserInfoRes{
  string id = 1;
  string name = 2;
  int32 age = 3;
}

message ReportUserBehaviorReq{
  string id = 1;
  string behavior = 2;
}

message ReportUserBehaviorRes{
  int32 retCode = 1;
  string retMsg = 2;
  int64 receivedCount = 3;
}

message SupplyUserChangeInfoReq{
  string id = 1;
}

message SupplyUserChangeInfoRes{
  string id = 1;
  string loc = 2;
  string temperature = 3;
}

service UserInfoService {
  rpc getUserInfo(GetUserInfoReq) returns (GetUserInfoRes);
  rpc reportUserBehavior(stream ReportUserBehaviorReq) returns (ReportUserBehaviorRes);
//服务端单向流
  rpc supplyUserChangeInfo(SupplyUserChangeInfoReq) returns (stream SupplyUserChangeInfoRes);
}
```

- **Java语言实现**

a、服务端实现接口对外服务

```java
package com.example.userinfoservice.impl;

import com.example.userinfoservice.UserInfoEntity;
import com.example.userinfoservice.UserInfoServiceGrpc;
import com.example.userinfoservice.entity.UserInfoData;
import com.example.userinfoservice.entity.UserLocTmp;
import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import io.grpc.stub.StreamObserver;
import org.apache.commons.lang3.StringUtils;

import java.util.HashMap;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

public class UserInfoServiceImpl extends UserInfoServiceGrpc.UserInfoServiceImplBase {

    private HashMap<String, UserInfoData> userinfoHashMap = Maps.newHashMap();

    private HashMap<String, List<UserLocTmp>> userLocTmpHashMap = Maps.newHashMap();

    @Override
    public void getUserInfo(UserInfoEntity.GetUserInfoReq request, StreamObserver<UserInfoEntity.GetUserInfoRes> responseObserver) {
        System.out.println("java服务收到请求");
        if (StringUtils.isEmpty(request.getId())){
            responseObserver.onError(new Exception("id不能为空"));
            responseObserver.onCompleted();
            return;
        }

        if (userinfoHashMap.isEmpty()){
            initUserInfoHashMap();
        }

        UserInfoData userInfoData = userinfoHashMap.get(request.getId());
        if (userInfoData == null){
            responseObserver.onError(new Exception("id不能为空"));
            responseObserver.onCompleted();
            return;
        }
        responseObserver.onNext(UserInfoEntity.GetUserInfoRes.newBuilder()
                .setId(userInfoData.getId())
                .setName(userInfoData.getName())
                .setAge(userInfoData.getAge())
                .build());
        responseObserver.onCompleted();
    }

    //客户端流
    @Override
    public StreamObserver<UserInfoEntity.ReportUserBehaviorReq> reportUserBehavior(StreamObserver<UserInfoEntity.ReportUserBehaviorRes> responseObserver) {
        AtomicInteger count = new AtomicInteger(0);
        UserInfoEntity.ReportUserBehaviorRes.Builder res = UserInfoEntity.ReportUserBehaviorRes.newBuilder()
                .setRetCode(0)
                .setRetMsg("收到~~");

        StreamObserver<UserInfoEntity.ReportUserBehaviorReq> reportUserBehaviorReqStreamObserver = new StreamObserver<UserInfoEntity.ReportUserBehaviorReq>() {
            @Override
            public void onNext(UserInfoEntity.ReportUserBehaviorReq value) {
                System.out.println(value.getId() + "\t" + value.getBehavior());
                count.incrementAndGet();
            }

            @Override
            public void onError(Throwable t) {
            }

            @Override
            public void onCompleted() {
                res.setReceivedCount(count.get());
                responseObserver.onNext(res.build());

                responseObserver.onCompleted();
            }
        };
        return reportUserBehaviorReqStreamObserver;
    }

    //服务端流
    @Override
    public void supplyUserChangeInfo(UserInfoEntity.SupplyUserChangeInfoReq request, StreamObserver<UserInfoEntity.SupplyUserChangeInfoRes> responseObserver) {
        String id = request.getId();
        if (userLocTmpHashMap.isEmpty()){
            initUserLocTmpHashMap();
        }
        List<UserLocTmp> userLocTmps = userLocTmpHashMap.get(id);
        for (UserLocTmp userLocTmp : userLocTmps) {
            responseObserver.onNext(UserInfoEntity.SupplyUserChangeInfoRes.newBuilder().setId(id)
            .setLoc(userLocTmp.getLoc())
            .setTemperature(userLocTmp.getTemperature())
            .build());
        }
        responseObserver.onCompleted();
    }

    private void initUserInfoHashMap() {
        userinfoHashMap.put("1",new UserInfoData("1","小美",24));
        userinfoHashMap.put("2",new UserInfoData("2","小强",28));
        userinfoHashMap.put("3",new UserInfoData("3","小刚",30));
    }
    private void initUserLocTmpHashMap() {
        UserLocTmp data = new UserLocTmp("上海", "24.5度");
        List<UserLocTmp> datas = Lists.newArrayList(data);
        for (int i = 0; i < 5; i++) {
            datas.add(data);
        }
        userLocTmpHashMap.put("1",datas);
        userLocTmpHashMap.put("2",datas);
        userLocTmpHashMap.put("3",datas);
    }

}
```

b、客户端接受服务端下发数据

```java
private static void SupplyUserChangeInfo(UserInfoServiceGrpc.UserInfoServiceStub asyncStub) {
        CountDownLatch finishLatch = new CountDownLatch(1);
        StreamObserver<UserInfoEntity.SupplyUserChangeInfoRes> resStreamObserver = new StreamObserver<UserInfoEntity.SupplyUserChangeInfoRes>() {
            @Override
            public void onNext(UserInfoEntity.SupplyUserChangeInfoRes value) {
                logger.info("id:"+value.getId()+"\tloc:"+value.getLoc()+"\ttemp:"+value.getTemperature());
            }
            @Override
            public void onError(Throwable t) {
            }
            @Override
            public void onCompleted() {
                logger.info("服务端流关闭");
                finishLatch.countDown();
            }
        };
        asyncStub.supplyUserChangeInfo(UserInfoEntity.SupplyUserChangeInfoReq.newBuilder().setId("1").build(),resStreamObserver);
        try {
            if (!finishLatch.await(10, TimeUnit.SECONDS)) {
                logger.warning("FAILED : Process orders cannot finish within 10 seconds");
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }
```

- **GO语言实现**

a、服务端支持

```go
func (server *StreamUserInfoServiceImpl) SupplyUserChangeInfo(req *pb.SupplyUserChangeInfoReq,
	res pb.UserInfoService_SupplyUserChangeInfoServer) error{
	//获取用户的数据
	if server.UserLocationTmp == nil{
		server.initUserLocationTmp()
	}

	id := req.Id
	if len(id) == 0 {
		return errors.New("id不能为空")
	}
	//获取该用户的历史移动轨迹
	locTmps := server.UserLocationTmp[id]
	for i, tmp := range locTmps {
		log.Printf("i:%d,data:%v",i,tmp)
		err := res.Send(&pb.SupplyUserChangeInfoRes{
			Id:id,
			Loc: tmp.Loc,
			Temperature: tmp.Temperature,
		})
		if err != nil {
			return fmt.Errorf("error sending message to stream : %v", err)
		}
	}
	return nil
}
```

b、客户端接收数据

```go
func SupplyUserChangeInfo(c pb.UserInfoServiceClient) {
	req:= &pb.SupplyUserChangeInfoReq{
		Id:"1",
	}
	res,err:=c.SupplyUserChangeInfo(context.Background(),req)
	if err != nil {
		log.Fatalf("err:%v",err)
		return
	}
	for{
		data,err:=res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		log.Printf("data:%v",data)
	}
}
```

## 3、客户端服务端双向流

在proto文件中创建一个服务端下发流数据的rpc方法 - **supplyUserChangeInfo**

```protobuf
syntax = "proto3";//标识 proto版本 建议使用proto3
package userinfoservice;//proto包名 避免命名冲突，也可以作为引入其他proto文件时使用
option java_package = "com.example.userinfoservice" ;//生成的类将带有此包名，不指定则使用package
option cc_generic_services = true;
option go_package = "./pb";
option java_outer_classname = "UserInfoEntity";

message GetUserInfoReq{
  string id = 1;
}
message GetUserInfoRes{
  string id = 1;
  string name = 2;
  int32 age = 3;
}

message ReportUserBehaviorReq{
  string id = 1;
  string behavior = 2;
}

message ReportUserBehaviorRes{
  int32 retCode = 1;
  string retMsg = 2;
  int64 receivedCount = 3;
}

message SupplyUserChangeInfoReq{
  string id = 1;
}

message SupplyUserChangeInfoRes{
  string id = 1;
  string loc = 2;
  string temperature = 3;
}

message ExchangeUserInfoReq{
  string id = 1;
  double lng = 2;
  double lat = 3;
}

message ExchangeUserInfoRes{
  string id = 1;
  string loc = 2;
}

service UserInfoService {
  rpc getUserInfo(GetUserInfoReq) returns (GetUserInfoRes);
  rpc reportUserBehavior(stream ReportUserBehaviorReq) returns (ReportUserBehaviorRes);
  rpc supplyUserChangeInfo(SupplyUserChangeInfoReq) returns (stream SupplyUserChangeInfoRes);
//双向流  
rpc exchangeUserInfo(stream ExchangeUserInfoReq) returns (stream ExchangeUserInfoRes);
}
```

- **Java语言实现**

a、服务端实现此方法

```java
//客户端服务端流
    @Override
    public StreamObserver<UserInfoEntity.ExchangeUserInfoReq> exchangeUserInfo(StreamObserver<UserInfoEntity.ExchangeUserInfoRes> responseObserver) {
        responseObserver.onNext(UserInfoEntity.ExchangeUserInfoRes.newBuilder()
                .setId("1")
                .setLoc("上海")
                .build());
        responseObserver.onNext(UserInfoEntity.ExchangeUserInfoRes.newBuilder()
                .setId("1")
                .setLoc("上海东站")
                .build());
        responseObserver.onNext(UserInfoEntity.ExchangeUserInfoRes.newBuilder()
                .setId("1")
                .setLoc("上海西站")
                .build());

        responseObserver.onNext(UserInfoEntity.ExchangeUserInfoRes.newBuilder()
                .setId("1")
                .setLoc("上海南站")
                .build());
        return new StreamObserver<UserInfoEntity.ExchangeUserInfoReq>() {
            @Override
            public void onNext(UserInfoEntity.ExchangeUserInfoReq value) {
                System.out.println("id:"+value.getId()+"\tlat:"+value.getLat()+"\tlng:"+value.getLng());
            }

            @Override
            public void onError(Throwable t) {

            }

            @Override
            public void onCompleted() {
                responseObserver.onCompleted();
            }
        };
    }
```

b、客户端接受流并上报流数据

```java
private static void ExchangeUserInfo(UserInfoServiceGrpc.UserInfoServiceStub asyncStub) {
        CountDownLatch finishLatch = new CountDownLatch(1);
        StreamObserver<UserInfoEntity.ExchangeUserInfoRes> resStreamObserver = new StreamObserver<UserInfoEntity.ExchangeUserInfoRes>() {
            @Override
            public void onNext(UserInfoEntity.ExchangeUserInfoRes value) {
                logger.info("res,id:"+value.getId()+"\tloc:"+value.getLoc());
            }

            @Override
            public void onError(Throwable t) {

            }

            @Override
            public void onCompleted() {
                logger.info("服务端流关闭 交换数据结束");
                finishLatch.countDown();
            }
        };
        StreamObserver<UserInfoEntity.ExchangeUserInfoReq> exchangeUserInfoReqStreamObserver = asyncStub.exchangeUserInfo(resStreamObserver);
        exchangeUserInfoReqStreamObserver.onNext(UserInfoEntity.ExchangeUserInfoReq.newBuilder().setId("1")
                .setLat(12.12)
                .setLng(23.23)
                .build());
        exchangeUserInfoReqStreamObserver.onNext(UserInfoEntity.ExchangeUserInfoReq.newBuilder().setId("1")
                .setLat(22.12)
                .setLng(33.23)
                .build());
        exchangeUserInfoReqStreamObserver.onNext(UserInfoEntity.ExchangeUserInfoReq.newBuilder().setId("1")
                .setLat(32.12)
                .setLng(43.23)
                .build());
        exchangeUserInfoReqStreamObserver.onNext(UserInfoEntity.ExchangeUserInfoReq.newBuilder().setId("1")
                .setLat(42.12)
                .setLng(53.23)
                .build());
        exchangeUserInfoReqStreamObserver.onCompleted();

        try {
            if (!finishLatch.await(10, TimeUnit.SECONDS)) {
                logger.warning("FAILED : Process orders cannot finish within 10 seconds");
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

    }
```

- **Go语言实现**

a、服务端实现

```go
func (server *StreamUserInfoServiceImpl) ExchangeUserInfo(ex pb.UserInfoService_ExchangeUserInfoServer) error {
	for  {
		recv, err := ex.Recv()
		if err == io.EOF{
			err := ex.Send(&pb.ExchangeUserInfoRes{
				Id:  "结束汇报",
				Loc: "结束汇报",
			})
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("用户Id:%s, lat:%v,lng:%v",recv.Id,recv.Lat,recv.Lng)
		err = ex.Send(&pb.ExchangeUserInfoRes{
			Id:  recv.Id,
			Loc: "收到，您当前的位置是北京",
		})
		if err != nil {
			return err
		}
	}
}
```

b、客户端实现

```go
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

	go listenExchangeUserInfo(info,lisChanClose)

	if err := info.CloseSend(); err != nil {
		log.Fatal(err)
	}

	 lisChanClose<-true

}
```