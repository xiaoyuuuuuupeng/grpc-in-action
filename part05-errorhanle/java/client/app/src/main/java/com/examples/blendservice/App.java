
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
