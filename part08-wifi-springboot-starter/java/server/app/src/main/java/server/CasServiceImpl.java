package server;

import com.example.casservice.CasServiceGrpc;
import com.example.casservice.CasServicePbEntity;
import io.grpc.stub.StreamObserver;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class CasServiceImpl extends CasServiceGrpc.CasServiceImplBase {
    @Override
    public void casLogin(CasServicePbEntity.CasLoginReq request, StreamObserver<CasServicePbEntity.CasLoginRes> responseObserver) {

        System.out.println(request.getUserName());
        responseObserver.onNext(CasServicePbEntity.CasLoginRes.newBuilder().setUid(UUID.randomUUID().toString()).build());
        responseObserver.onCompleted();

    }
}
