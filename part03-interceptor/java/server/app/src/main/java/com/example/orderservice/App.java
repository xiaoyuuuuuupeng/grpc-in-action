
package com.example.orderservice;

import com.example.orderservice.impl.OrderServiceImpl;
import com.example.orderservice.interceptors.ParamsInterceptor;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.protobuf.services.ProtoReflectionService;

import java.io.IOException;

public class App {

    public static void main(String[] args) throws IOException, InterruptedException {

        Server server = ServerBuilder.forPort(10082).intercept(new ParamsInterceptor()).addService(new OrderServiceImpl())
                .addService(ProtoReflectionService.newInstance())
                .build();
        server.start();
        server.awaitTermination();
    }
}
