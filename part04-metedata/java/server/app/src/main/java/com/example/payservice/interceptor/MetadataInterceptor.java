package com.example.payservice.interceptor;

import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;

public class MetadataInterceptor implements ServerInterceptor {
    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(ServerCall<ReqT, RespT> call, Metadata headers, ServerCallHandler<ReqT, RespT> next) {
        String token = headers.get(Metadata.Key.of("token", Metadata.ASCII_STRING_MARSHALLER));
        System.out.println("token:"+token);
        return next.startCall(call, headers);
    }
}
