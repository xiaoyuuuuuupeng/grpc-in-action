
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
