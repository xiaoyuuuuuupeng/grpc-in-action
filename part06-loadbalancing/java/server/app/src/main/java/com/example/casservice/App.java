
package com.example.casservice;

import com.example.casservice.impl.CasServiceImpl;
import com.google.common.collect.Lists;
import io.grpc.Server;
import io.grpc.ServerBuilder;

import java.io.IOException;
import java.util.List;

public class App {


    public static void main(String[] args) throws IOException, InterruptedException {
        List<Integer> ports = Lists.newArrayList();

        List<Server> servers = Lists.newArrayList();
        ports.add(10084);
        ports.add(10085);

        for (Integer port : ports) {
            System.out.println(port + "启动成功...");
            CasServiceImpl casService = new CasServiceImpl(port);
            Server server = ServerBuilder.forPort(port).addService(casService)
                    .build();
            server.start();
            servers.add(server);
        }

        for (Server server : servers) {
            server.awaitTermination();
        }

    }
}
