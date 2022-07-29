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
