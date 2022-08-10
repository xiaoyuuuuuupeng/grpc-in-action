
package com.exapmles.blendservice;

import com.exapmles.blendservice.impl.BlendServiceImpl;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.protobuf.services.ProtoReflectionService;

import java.io.IOException;

public class App {


    public static void main(String[] args) throws InterruptedException, IOException {

        Server server = ServerBuilder.forPort(10085)
                .addService(new BlendServiceImpl()).addService(ProtoReflectionService.newInstance()).build();
        server.start();
        server.awaitTermination();


    }
}
