package com.examples.easyservice;

import com.examples.easyservice.resolver.NacosNameResolverProvider;
import io.grpc.Attributes;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

import java.net.URI;

public class App {

    public static void main(String[] args) {
        URI uri = URI.create("http://localhost:8848");


        ManagedChannel managedChannel = ManagedChannelBuilder.forTarget("nacos://"+EasyServiceGrpc.getServiceDescriptor().getName()).
                defaultLoadBalancingPolicy("round_robin").
                nameResolverFactory(new NacosNameResolverProvider(uri, Attributes.newBuilder().build())).usePlaintext().build();

        EasyServiceGrpc.EasyServiceBlockingStub easyServiceBlockingStub = EasyServiceGrpc.newBlockingStub(managedChannel);

        for (int i = 0; i < 10; i++) {
            EasyServicePbEntity.Response response = easyServiceBlockingStub.sendRequest(EasyServicePbEntity.Request.newBuilder().setPing("ping from pengxiaoyu").build());
            System.out.println(response.getPong());
        }

    }
}
