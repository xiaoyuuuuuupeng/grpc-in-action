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
        System.out.println("getUserInfo");
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
        System.out.println("reportUserBehavior");

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
        System.out.println("supplyUserChangeInfo");
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
