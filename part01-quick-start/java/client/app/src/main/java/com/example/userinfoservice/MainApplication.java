package com.example.userinfoservice;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

import java.util.logging.Logger;

public class MainApplication {

    private static final Logger logger = Logger.getLogger(MainApplication.class.getName());

    public static void main(String[] args) throws InterruptedException {
        ManagedChannel channel = ManagedChannelBuilder.forAddress("localhost", 10080)
                .usePlaintext()
                .build();

        UserInfoServiceGrpc.UserInfoServiceBlockingStub stub =
                UserInfoServiceGrpc.newBlockingStub(channel);

        Userinfo.GetUserInfoRes res = stub.getUserInfo(
                Userinfo.GetUserInfoReq.newBuilder()
                        .setId("1")
                        .build());
        logger.info("用户id: " + res.getId() + "\t用户姓名:"+ res.getName()
        +"\t用户年龄:"+ res.getAge());
        channel.shutdown();
    }
}
