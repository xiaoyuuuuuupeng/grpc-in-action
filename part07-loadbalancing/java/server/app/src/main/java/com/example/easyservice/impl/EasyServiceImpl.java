package com.example.easyservice.impl;

import com.example.easyservice.EasyServiceGrpc;
import com.example.easyservice.EasyServicePbEntity;
import io.grpc.stub.StreamObserver;

public class EasyServiceImpl extends EasyServiceGrpc.EasyServiceImplBase {
    @Override
    public void sendRequest(EasyServicePbEntity.Request request, StreamObserver<EasyServicePbEntity.Response> responseObserver) {
        String ping = request.getPing();
        responseObserver.onNext(EasyServicePbEntity.Response.newBuilder()
                .setPong("get pong,from :"+ping).build());
        responseObserver.onCompleted();
    }
}
