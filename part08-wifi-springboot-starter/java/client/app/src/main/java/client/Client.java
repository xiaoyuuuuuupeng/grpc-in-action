package client;

import com.example.casservice.CasServiceGrpc;
import com.example.casservice.CasServicePbEntity;
import org.springframework.stereotype.Component;
import io.grpc.Channel;
import javax.annotation.Resource;

@Component
public class Client {
    @Resource
    private Channel apGreeterChannel; // bean名为: ${target名}Channel


    public String casLogin(String name) {
        CasServicePbEntity.CasLoginReq request = CasServicePbEntity.CasLoginReq.newBuilder().setUserName(name).build();
        CasServicePbEntity.CasLoginRes reply = CasServiceGrpc.newBlockingStub(apGreeterChannel).casLogin(request);
        return reply.getUid();
    }
}