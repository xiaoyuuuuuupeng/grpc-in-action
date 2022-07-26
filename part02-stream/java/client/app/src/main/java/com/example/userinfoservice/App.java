package com.example.userinfoservice;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.stub.StreamObserver;
import org.checkerframework.checker.units.qual.C;

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
        UserInfoServiceGrpc.UserInfoServiceStub asyncStub = UserInfoServiceGrpc.newStub(channel);
//        GetUserInfo(stub);
//        //客户端流

//        ReportUserBehavior(asyncStub);
//        //服务端流接收
//        SupplyUserChangeInfo(asyncStub);
        //客户端服务端双向接受
        ExchangeUserInfo(asyncStub);
        channel.shutdown();
    }

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

    private static void ReportUserBehavior(UserInfoServiceGrpc.UserInfoServiceStub stub) {
        CountDownLatch finishLatch = new CountDownLatch(1);
        StreamObserver<UserInfoEntity.ReportUserBehaviorRes> responseObserver = new StreamObserver<UserInfoEntity.ReportUserBehaviorRes>() {
            @Override
            public void onNext(UserInfoEntity.ReportUserBehaviorRes value) {
                logger.info("Report User Behavior RetCode: " + value.getRetCode() + " RetMsg: " + value.getRetMsg() + "count:" + value.getReceivedCount());
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
