package com.example.userinfoservice.impl;

import com.example.userinfoservice.UserInfoEntity;
import com.example.userinfoservice.UserInfoServiceGrpc;
import com.example.userinfoservice.entity.UserInfoData;
import com.google.common.collect.Maps;
import io.grpc.stub.StreamObserver;
import org.apache.commons.lang3.StringUtils;

import java.util.HashMap;

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

    private void initUserInfoHashMap() {
        userinfoHashMap.put("1",new UserInfoData("1","小美",24));
        userinfoHashMap.put("2",new UserInfoData("2","小强",28));
        userinfoHashMap.put("3",new UserInfoData("3","小刚",30));
    }
}
