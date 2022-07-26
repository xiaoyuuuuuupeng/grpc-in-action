package com.example.orderservice.interceptors;

import io.grpc.*;

import java.util.logging.Logger;

public class ParamsInterceptor implements ServerInterceptor {
    Logger logger = Logger.getLogger(this.getClass().getName());
    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(ServerCall<ReqT, RespT> call, Metadata headers, ServerCallHandler<ReqT, RespT> next) {
        logger.info("======= [Server Interceptor] : Remote Method Invoked - " + call.getMethodDescriptor().getFullMethodName());
        ServerCall<ReqT, RespT> serverCall = new OrderMsgServerCall<>(call);
        return new OrderMsgServerCallListener<>(next.startCall(serverCall, headers));
    }
}
