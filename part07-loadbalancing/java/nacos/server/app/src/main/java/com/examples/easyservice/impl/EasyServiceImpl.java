package com.examples.easyservice.impl;

import com.examples.easyservice.EasyServiceGrpc;
import com.examples.easyservice.EasyServicePbEntity;
import com.examples.easyservice.utils.NetUtils;
import io.grpc.stub.StreamObserver;

public class EasyServiceImpl extends EasyServiceGrpc.EasyServiceImplBase {
    int port = 0;

    public EasyServiceImpl(int port) {
        this.port = port;
    }

    @Override
    public void sendRequest(EasyServicePbEntity.Request request, StreamObserver<EasyServicePbEntity.Response> responseObserver) {
        String ping = request.getPing();
        System.out.println("ping======>"+ping);
        responseObserver.onNext(EasyServicePbEntity.Response.newBuilder().setPong("pong from" + NetUtils.getLocalHost()+ ":"+port).build());
        responseObserver.onCompleted();
    }
}
