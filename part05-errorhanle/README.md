# 五、错误处理

在业务开发中，大家是会使用http通用的错误码，还是统一返回HttpCode-200，转而使用业务自己定义的错误码呢？这个问题可能见仁见智，大家都有一些自己的理解和最佳实践。这节我们来看下gRPC中如何对错误进行处理的。

gRPC提供了一组定义良好的专用状态码，举例如下：

## 一、一般的错误

| 案例 | 状态码 |
| --- | --- |
| 客户端应用程序取消了请求 | GRPC_STATUS_CANCELLED |
| 截止日期在服务器返回状态之前已过期 | GRPC_STATUS_DEADLINE_EXCEEDED |
| 在服务器上找不到方法 | GRPC_STATUS_UNIMPLEMENTED |
| 服务器关闭 | GRPC_STATUS_UNAVAILABLE |
| 服务器抛出异常（或者除了返回状态码来终止 RPC 之外做了其他事情） | GRPC_STATUS_UNKNOWN |

## 二、网络故障

| 案例 | 状态码 |
| --- | --- |
| 在截止日期到期之前不会传输任何数据。也适用于在截止日期到期之前传输了一些数据并且没有检测到其他故障的情况 | GRPC_STATUS_DEADLINE_EXCEEDED |
| 连接中断前传输的一些数据（例如，请求元数据已写入 TCP 连接） | GRPC_STATUS_UNAVAILABLE |

## 三、协议错误

| 案例 | 状态码 |
| --- | --- |
| 无法解压但支持压缩算法 | GRPC_STATUS_INTERNAL |
| 服务器不支持客户端使用的压缩机制 | GRPC_STATUS_UNIMPLEMENTED |
| 达到流量控制资源限制 | GRPC_STATUS_RESOURCE_EXHAUSTED |
| 流量控制协议违规 | GRPC_STATUS_INTERNAL |
| 解析返回状态时出错 | GRPC_STATUS_UNKNOWN |
| 未经身份验证：凭据无法获取元数据 | GRPC_STATUS_UNAUTHENTICATED |
| 权限元数据中的主机集无效 | GRPC_STATUS_UNAUTHENTICATED |
| 解析响应协议缓冲区时出错 | GRPC_STATUS_INTERNAL |
| 解析请求协议缓冲区时出错 | GRPC_STATUS_INTERNAL |

## 四、代码实现

这里模拟前端请求推荐列表页的接口，当用户没有登录是UserId=-1，这时候服务端因为没有用户特征所以会拒绝下发数据，转而返回一个gRPC特定错误。现在来看下，Java和go是分别怎么样实现的吧。

### 1、Go 服务端

在返回错误码的时候，还可以返回绑定的错误详情

这里使用了gRPC自带的包

```go
import (
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-in-action/part05-errorhanle/go/server/pb"
	"log"
)
```

```go
func (*BlendService) Blend(context context.Context, req *pb.BlendReq) (*pb.BlendRes, error)  {

	if req.UserId == "-1" {
		log.Printf("UserId is invalid! -> Received User ID %s", req.UserId)
		errorStatus := status.New(codes.InvalidArgument, "Invalid information received")
		ds, err := errorStatus.WithDetails(
			&errdetails.BadRequest_FieldViolation{
				Field:"ID",
				Description: fmt.Sprintf("UserId received is not valid :%s ", req.UserId),
			},
		)
		if err != nil {
			return nil, errorStatus.Err()
		}
		return nil, ds.Err()
	}else {
		nids :=[]string{"1","2","3"}

		return &pb.BlendRes{
			Nid:nids,
		},nil
	}
}
```

### 2、Go客户端

客户端的行为就是拿到服务端返回的错误码和错误详情，进行对应的处理即可

```go
package main

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"grpc-in-action/part05-errorhanle/go/client/pb"
	"log"
)

func main() {

	dial, err := grpc.Dial("localhost:10085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewBlendServiceClient(dial)
	blend, err := client.Blend(context.Background(), &pb.BlendReq{
		UserId: "-1",
	})
	if err != nil {
		log.Printf("err is :%+v",err)
		errorCode := status.Code(err)
		if errorCode == codes.InvalidArgument {
			log.Printf("Invalid Argument Error : %s", errorCode)
			errorStatus := status.Convert(err)
			for _, d := range errorStatus.Details() {
				switch info := d.(type) {
				case *errdetails.BadRequest_FieldViolation:
					log.Printf("Request Field Invalid: %s", info)
				default:
					log.Printf("Unexpected error type: %s", info)
				}
			}
		} else {
			log.Printf("Unhandled error : %s ", errorCode)
		}
	}else {
		log.Printf("blend:%v",blend.Nid)
	}

}
```

### 3、Java服务端

java是通过异常的方式对外发送一个错误信息

```java
public class BlendServiceImpl extends BlendServiceGrpc.BlendServiceImplBase {

    private static final Logger logger = Logger.getLogger(BlendServiceImpl.class.getName());
    @Override
    public void blend(BlendServicePbEntity.blendReq request, StreamObserver<BlendServicePbEntity.blendRes> responseObserver) {
        if ("-1".equals(request.getUserId())){
            logger.info("Invalid Order ID: " + request.getUserId());
            responseObserver.onError(Status.INVALID_ARGUMENT.withDescription("Invalid UserId received.").asException());
        }

        responseObserver.onNext(BlendServicePbEntity.blendRes.newBuilder().addNid("java_1")
                .addNid("java_2")
                .addNid("java_3")
                .build());
        responseObserver.onCompleted();
    }
}
```

### 4、Java客户端

接收当然也是使用处理异常的方式

```java
package com.examples.blendservice;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.StatusRuntimeException;

import java.util.logging.Logger;

public class App {
    private static final Logger logger = Logger.getLogger(App.class.getName());

    public static void main(String[] args) {

        ManagedChannel managedChannel = ManagedChannelBuilder.forAddress("localhost",10085)
                .usePlaintext().build();
        try{
            BlendServicePbEntity.blendRes blend = BlendServiceGrpc.newBlockingStub(managedChannel).blend(BlendServicePbEntity.blendReq.newBuilder()
                    .setUserId("-1")
                    .build());
            for (String nid : blend.getNidList()) {
                System.out.println(nid);
            }
        }catch (StatusRuntimeException e) {
            logger.info(" Error Received - Error Code : " + e.getStatus().getCode());
            logger.info("Error details -> " + e.getMessage());
        }

    }
}
```

- 参考文档

[https://grpc.io/docs/guides/error/#protocol-errors](https://grpc.io/docs/guides/error/#protocol-errors)

[https://github.com/avinassh/grpc-errors](https://github.com/avinassh/grpc-errors)