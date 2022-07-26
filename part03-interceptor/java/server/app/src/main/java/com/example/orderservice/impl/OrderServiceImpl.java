package com.example.orderservice.impl;

import com.example.orderservice.OrderPbEntity;
import com.example.orderservice.OrderServiceGrpc;
import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import io.grpc.stub.StreamObserver;

import java.util.HashMap;
import java.util.List;

public class OrderServiceImpl extends OrderServiceGrpc.OrderServiceImplBase {

    private HashMap<String, List<String>> data = Maps.newHashMap();

    @Override
    public void getOrderList(OrderPbEntity.GetOrderListReq request, StreamObserver<OrderPbEntity.GetOrderListRes> responseObserver) {
        if (data.isEmpty()){
            initData();
        }
        String userId = request.getUserId();
        List<String> strings = data.get(userId);
        responseObserver.onNext(OrderPbEntity.GetOrderListRes.newBuilder().addAllOrderIds(strings).build());
        responseObserver.onCompleted();
    }

    private void initData() {
        List<String> data1 = Lists.newArrayList();
        data1.add("order_1");
        data1.add("order_2");
        data1.add("order_3");
        data1.add("order_4");
        data.put("1",data1);
        data.put("2",data1);


    }
}
