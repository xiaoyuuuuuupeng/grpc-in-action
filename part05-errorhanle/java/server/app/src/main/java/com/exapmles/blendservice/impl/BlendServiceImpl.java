package com.exapmles.blendservice.impl;


import com.examples.blendservice.BlendServiceGrpc;
import com.examples.blendservice.BlendServicePbEntity;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;

import java.util.logging.Logger;

public class BlendServiceImpl extends BlendServiceGrpc.BlendServiceImplBase {

    private static final Logger logger = Logger.getLogger(BlendServiceImpl.class.getName());
    @Override
    public void blend(BlendServicePbEntity.blendReq request, StreamObserver<BlendServicePbEntity.blendRes> responseObserver) {
        if ("-1".equals(request.getUserId())){
            logger.info("Invalid Order ID: " + request.getUserId());
            responseObserver.onError(Status.INVALID_ARGUMENT.withDescription("Invalid UserId received.").asException());
        }

        responseObserver.onNext(BlendServicePbEntity.blendRes.newBuilder().addNid("java_1")
                .addNid("java_2")
                .addNid("java_3")
                .build());
        responseObserver.onCompleted();
    }
}
