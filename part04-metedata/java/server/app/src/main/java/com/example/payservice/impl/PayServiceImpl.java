package com.example.payservice.impl;

import com.example.payservice.PayServiceGrpc;
import com.example.payservice.PayServicePbEntity;
import io.grpc.stub.StreamObserver;

public class PayServiceImpl extends PayServiceGrpc.PayServiceImplBase {


    @Override
    public void payOrder(PayServicePbEntity.payOrderReq request, StreamObserver<PayServicePbEntity.payOrderRes> responseObserver) {

        System.out.println("payOrder,order_Id:"+request.getOrderId());

        responseObserver.onNext(PayServicePbEntity.payOrderRes.newBuilder().setRetCode(0)
        .setRetMsg("支付成功")
        .build());
        responseObserver.onCompleted();
    }
}
