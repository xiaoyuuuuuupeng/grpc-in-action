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
