# 六、负载均衡（一）

生产环境中，往往生产者和消费者都是多个，消费者需要从服务列表中挑选一个生产者进行消费，进而有了微服务环境中的一个特性负载均衡。RPC框架基本都支持负载均衡，duboo这方面做的很优秀，框架内部就支持了多种负载均衡策略。接下来，我们来介绍下负载均衡在grpc里的设计和思考。相信大家肯定在其中有所收获。

## 一、gRPC自带负载均衡算法

### 1、定义proto接口

```protobuf
syntax = "proto3";//标识 proto版本 建议使用proto3
package casservice;//proto包名 避免命名冲突，也可以作为引入其他proto文件时使用
option java_package = "com.example.casservice" ;//生成的类将带有此包名，不指定则使用package
option cc_generic_services = true;
option go_package = "./pb";
option java_outer_classname = "CasServicePbEntity";

message CasLoginReq{
  string userName = 1;
  string password = 2;
}
message CasLoginRes{
   string uid = 1;
   string sex = 2;
}

service CasService {
  rpc casLogin(CasLoginReq) returns (CasLoginRes);
}
```

### 2、Java语言实现

- 服务端实现

```java
package com.example.casservice.impl;

import com.example.casservice.CasServiceGrpc;
import com.example.casservice.CasServicePbEntity;
import io.grpc.stub.StreamObserver;

import java.util.UUID;

public class CasServiceImpl extends CasServiceGrpc.CasServiceImplBase {

    public CasServiceImpl(Integer port) {
        this.port = port;
    }

    private Integer port = 0 ;
    @Override
    public void casLogin(CasServicePbEntity.CasLoginReq request, StreamObserver<CasServicePbEntity.CasLoginRes> responseObserver) {

        System.out.println("port:"+port+"登录请求，用户名："+request.getUserName()+"\t密码："+request.getPassword());
        responseObserver.onNext(CasServicePbEntity.CasLoginRes.newBuilder()
                .setUid("port:"+port+"\t"+UUID.randomUUID().toString().replaceAll("-",""))
                .setSex("男").build());
        responseObserver.onCompleted();

    }
}
```

- 启动多个服务端实例

```java
package com.example.casservice;

import com.example.casservice.impl.CasServiceImpl;
import com.google.common.collect.Lists;
import io.grpc.Server;
import io.grpc.ServerBuilder;

import java.io.IOException;
import java.util.List;

public class App {
    
    public static void main(String[] args) throws IOException, InterruptedException {
        List<Integer> ports = Lists.newArrayList();

        List<Server> servers = Lists.newArrayList();
        ports.add(10084);
        ports.add(10085);

        for (Integer port : ports) {
            System.out.println(port + "启动成功...");
            CasServiceImpl casService = new CasServiceImpl(port);
            Server server = ServerBuilder.forPort(port).addService(casService)
                    .build();
            server.start();
            servers.add(server);
        }

        for (Server server : servers) {
            server.awaitTermination();
        }

    }
}
```

- 客户端需要实现一个自定的namereslover来将请求的域名解析为服务实例ip

```java
package com.example.casservice.resolve;

import com.google.common.collect.Lists;
import io.grpc.EquivalentAddressGroup;
import io.grpc.NameResolver;

import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

import io.grpc.Attributes;

public class CustomNameResolver extends NameResolver {

    private final String path;

    private Listener2 listener;

    public CustomNameResolver(String path) {
        this.path = path;
    }

    @Override
    public String getServiceAuthority() {
        return this.path;
    }

    @Override
    public void shutdown() {
    }

    @Override
    public void start(Listener2 listener) {
        this.listener = listener;
        this.resolve();
    }

    private void resolve() {
        List<SocketAddress> socketAddressList = getAddressList(path).stream()
                .map(this::toSocketAddress)
                .collect(Collectors.toList());

        List<EquivalentAddressGroup> equivalentAddressGroups = Lists.newArrayList();
        for (SocketAddress socketAddress : socketAddressList) {
            EquivalentAddressGroup equivalentAddressGroup = new EquivalentAddressGroup(socketAddress);
            equivalentAddressGroups.add(equivalentAddressGroup);
        }
        ResolutionResult resolutionResult = ResolutionResult.newBuilder()
                .setAddresses(equivalentAddressGroups)
                .setAttributes(Attributes.EMPTY)
                .build();

        this.listener.onResult(resolutionResult);

    }
    private SocketAddress toSocketAddress(InetSocketAddress address) {
        return new InetSocketAddress(address.getHostName(), address.getPort());
    }

    private List<InetSocketAddress> getAddressList(String path) {
        System.out.println("path:"+path);
        InetSocketAddress inetSocketAddress = new InetSocketAddress("localhost", 10084);
        InetSocketAddress inetSocketAddress2 = new InetSocketAddress("localhost", 10085);
        return Arrays.asList(inetSocketAddress, inetSocketAddress2);
    }

}
```

- 新规范需要使用provider注册nameresolver

```java
package com.example.casservice.resolve;

import io.grpc.NameResolver;
import io.grpc.NameResolverProvider;
import java.net.URI;

public class CustomNameResolverProvider extends NameResolverProvider {

    public CustomNameResolverProvider() {

    }
    @Override
    public NameResolver newNameResolver(URI targetUri, NameResolver.Args args) {
        return new CustomNameResolver(targetUri.getPath());
    }

    @Override
    protected boolean isAvailable() {
        return true;
    }

    @Override
    protected int priority() {
        return 10;
    }

    @Override
    public String getDefaultScheme() {
        return "example";
    }
}
```

- `使用 **round_robin、**`**`pick_first` 轮询的方式请求 发起二十次请求**

```java
package com.example.casservice;

import com.example.casservice.resolve.CustomNameResolverProvider;
import io.grpc.*;

public class App {

    public static void main(String[] args) {

        NameResolverRegistry.getDefaultRegistry().register(new CustomNameResolverProvider());
        ManagedChannel channel = ManagedChannelBuilder.forTarget("example:///lb.example.grpc.io")
                .defaultLoadBalancingPolicy("pick_first")
                .usePlaintext()
                .build();
        CasServiceGrpc.CasServiceBlockingStub casServiceBlockingStub = CasServiceGrpc.newBlockingStub(channel);
        System.out.println("==== Calling with pick_first ====");

        for (int i = 0; i < 20; i++) {

            CasServicePbEntity.CasLoginRes casLoginRes =casServiceBlockingStub.casLogin(CasServicePbEntity.CasLoginReq.newBuilder()
                    .setUserName("java_test")
                    .setPassword("test_password").build());

            System.out.println("casLoginRes,uid:"+casLoginRes.getUid()+"性别:\t"+casLoginRes.getSex());
        }

        System.out.println("==== Calling  with round_robin ====");
        ManagedChannel channel2 = ManagedChannelBuilder.forTarget("example:///lb.example.grpc.io")
                .defaultLoadBalancingPolicy("round_robin")
                .usePlaintext()
                .build();
        CasServiceGrpc.CasServiceBlockingStub casServiceBlockingStub2 = CasServiceGrpc.newBlockingStub(channel2);
        for (int i = 0; i < 20; i++) {

            CasServicePbEntity.CasLoginRes casLoginRes =casServiceBlockingStub2.casLogin(CasServicePbEntity.CasLoginReq.newBuilder()
                    .setUserName("java_test2")
                    .setPassword("test_password2").build());

            System.out.println("casLoginRes,uid:"+casLoginRes.getUid()+"性别:\t"+casLoginRes.getSex());
        }

    }
}
```

```bash
==== Calling with pick_first ====
path:/lb.example.grpc.io
casLoginRes,uid:port:10084	441f1ce54214475cb366982e3e0fd726性别:	男
casLoginRes,uid:port:10084	a4cfed898db64706a18507983d8589d6性别:	男
casLoginRes,uid:port:10084	a04fb2b71944440b8ed00f02bdc3b15f性别:	男
casLoginRes,uid:port:10084	2b2668b2f1ad472ab35729708f0959d3性别:	男
casLoginRes,uid:port:10084	1fa78e4526c24daf9c6d9779d1626ca9性别:	男
casLoginRes,uid:port:10084	826faf498afc4e1ba65db18a1eb0ccf4性别:	男
casLoginRes,uid:port:10084	c9c992e0b4954b9a9f1e94022c8a65ca性别:	男
casLoginRes,uid:port:10084	08b52ac7e444450c9aaad323d255b168性别:	男
casLoginRes,uid:port:10084	9cde040761cb4de993f3fbc95e9299e3性别:	男
casLoginRes,uid:port:10084	d8aa79dbcd094c0f918b219ba3cbb7a2性别:	男
casLoginRes,uid:port:10084	78575fcb8e2f40cab1a3085e7d609c1f性别:	男
casLoginRes,uid:port:10084	c1863d7b914045d6b339a5be8770b215性别:	男
casLoginRes,uid:port:10084	80bd8a2a1a444ade901c849a7e338e63性别:	男
casLoginRes,uid:port:10084	5e37a848792a41dab612b13089677771性别:	男
casLoginRes,uid:port:10084	1674a08cdcc64b4882805038485742d5性别:	男
casLoginRes,uid:port:10084	192055f9ebfd4a1db03679cc875560c6性别:	男
casLoginRes,uid:port:10084	63d611cab7884ff7b3e5e21073fc6164性别:	男
casLoginRes,uid:port:10084	d2d7ca2e4122428a8d134160e47938e9性别:	男
casLoginRes,uid:port:10084	6e924aae0c3e44348dce641c06652440性别:	男
casLoginRes,uid:port:10084	91a9dd7f159e49a0948ce30c4fb1ecf2性别:	男
==== Calling  with round_robin ====
path:/lb.example.grpc.io
casLoginRes,uid:port:10085	f3d4133b24ac4ab1ae2afac4078e46e7性别:	男
casLoginRes,uid:port:10084	075838a14fdc4626a07892d3701c30ce性别:	男
casLoginRes,uid:port:10085	a627d8c99a2b4b1185e5414565142611性别:	男
casLoginRes,uid:port:10084	58cca91d354b4688807674278cda539e性别:	男
casLoginRes,uid:port:10085	3be7bcd4a43247b4b889feb4daf8724f性别:	男
casLoginRes,uid:port:10084	b7ff85a74db44736804234bc20e1feb8性别:	男
casLoginRes,uid:port:10085	ba0028287c0c4a62b2854c25c1344120性别:	男
casLoginRes,uid:port:10084	f6a045814e2149e9833ae644d563c87b性别:	男
casLoginRes,uid:port:10085	0596a4fc0b644d5a9e734ec9262abf7f性别:	男
casLoginRes,uid:port:10084	574ff946c98547e7a23d30dd9905d9fa性别:	男
casLoginRes,uid:port:10085	b7d6a68b5e344310ae133b11a2b6635e性别:	男
casLoginRes,uid:port:10084	9ab058a644fd45eb9ea90fbb7df0b575性别:	男
casLoginRes,uid:port:10085	1153074d7626465d96ef041dafa55bb3性别:	男
casLoginRes,uid:port:10084	b00ddaec91cd40a1abe9cd83a39bd5ef性别:	男
casLoginRes,uid:port:10085	840417637b264cefb34b284ce55daa65性别:	男
casLoginRes,uid:port:10084	d033dba7c8fb49cdacbe4ae7bd91c714性别:	男
casLoginRes,uid:port:10085	7ccdd8f59a4f448a91df00dff15ac258性别:	男
casLoginRes,uid:port:10084	e6042be3517e47f9a95257507d1aa66c性别:	男
casLoginRes,uid:port:10085	8bf6183990c2449cb7e58f11534e4211性别:	男
casLoginRes,uid:port:10084	54788c7324ce45eb88339c14f03c9be2性别:	男
```

### 2、Go语言客户端实现

```go
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"grpc-in-action/part06-loadbalancing/go/client/pb"
	"log"
)
const (
	exampleScheme      = "example"
	exampleServiceName = "lb.example.grpc.io"
)

var addrs = []string{"localhost:10084", "localhost:10085"}

func main() {
	pickfirstConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName), // "example:///lb.example.grpc.io"
		// grpc.WithBalancerName("pick_first"), // "pick_first" is the default, so this DialOption is not necessary.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer pickfirstConn.Close()

	log.Println("==== Calling with pick_first ====")
	makeRPCs(pickfirstConn, 10)

	// Make another ClientConn with round_robin policy.
	roundrobinConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName), // // "example:///lb.example.grpc.io"
		//grpc.WithBalancer(("round_robin"), // This sets the initial balancing policy.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer roundrobinConn.Close()

	log.Println("==== Calling  with round_robin ====")
	makeRPCs(roundrobinConn, 10)
}

func makeRPCs(conn *grpc.ClientConn, count int) {
	client := pb.NewCasServiceClient(conn)
	req:= &pb.CasLoginReq{
		UserName: "go_cas_req"}
	for i := 0; i <count ; i++ {
		res, err := client.CasLogin(context.Background(), req)
		if err != nil {
			log.Printf("err:%v",err)
		}
		log.Printf("res:%v",res)
	}
}

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			"/lb.example.grpc.io":{"localhost:10084", "localhost:10085"}, // "lb.example.grpc.io": "localhost:50051", "localhost:50052"
		},
	}
	r.start()
	return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return "example" } // "example"

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.URL.Path]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*exampleResolver) Close()                                 {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
```

## 二、小结

gRPC默认的负载均衡方式只有三种：****`pick_first`、`round_robin`、`grpclb` （**但这种方式已经过时了，被xds所替代）

gRPC提供的前两种方式，在复杂的生产环境中确实显得有些单薄，那作为一款优秀的框架，它是如何为我们提供解决方案的呢、

- 相关连接

[https://github.com/grpc/grpc/blob/master/doc/load-balancing.md](https://github.com/grpc/grpc/blob/master/doc/load-balancing.md)

[https://github.com/grpc/grpc/blob/master/doc/naming.md](https://github.com/grpc/grpc/blob/master/doc/naming.md)

[https://grpc.io/blog/grpc-load-balancing/](https://grpc.io/blog/grpc-load-balancing/)