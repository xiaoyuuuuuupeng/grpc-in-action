package com.example.orderservice.interceptors;

import com.example.orderservice.OrderPbEntity;
import io.grpc.ForwardingServerCallListener;
import io.grpc.ServerCall;

import java.util.logging.Logger;

public class OrderMsgServerCallListener<R> extends ForwardingServerCallListener<R> {
    private static final Logger logger = Logger.getLogger(OrderMsgServerCallListener.class.getName());

    private final ServerCall.Listener<R> delegate;

    OrderMsgServerCallListener(ServerCall.Listener<R> delegate) {
        this.delegate = delegate;
    }

    @Override
    protected ServerCall.Listener<R> delegate() {
        return delegate;
    }

    @Override
    public void onMessage(R message) {
        if(message instanceof OrderPbEntity.GetOrderListReq){
            OrderPbEntity.GetOrderListReq req = (OrderPbEntity.GetOrderListReq)message;
            String userId = req.getUserId();
            logger.info("userId:"+userId);
        }
        logger.info("Message Received from Client -> Service " + message);
        super.onMessage(message);
    }
}
