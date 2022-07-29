package com.example.casservice.resolve;

import io.grpc.EquivalentAddressGroup;
import io.grpc.NameResolver;

import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

public class CustomNameResolver extends NameResolver {
    private final String authority;

    private Listener2 listener;


    public CustomNameResolver(String authority) {
        this.authority = authority;
    }

    @Override
    public String getServiceAuthority() {
        return this.authority;
    }

    @Override
    public void shutdown() {
    }

    @Override
    public void start(Listener2 listener) {
        this.listener = listener;
        this.resolve();
    }

    @Override
    public void refresh() {
        this.resolve();
    }
    private void resolve() {
        List<SocketAddress> socketAddressList = getAddressList(authority).stream()
                .map(this::toSocketAddress)
                .collect(Collectors.toList());
        EquivalentAddressGroup equivalentAddressGroup = new EquivalentAddressGroup(socketAddressList);

        ResolutionResult resolutionResult = ResolutionResult.newBuilder()
                .setAddresses(Arrays.asList(equivalentAddressGroup))
                .build();

        this.listener.onResult(resolutionResult);

    }
    private SocketAddress toSocketAddress(InetSocketAddress address) {
        return new InetSocketAddress(address.getHostName(), address.getPort());
    }


    private List<InetSocketAddress> getAddressList(String authority) {
        InetSocketAddress inetSocketAddress = new InetSocketAddress("localhost", 10084);
        InetSocketAddress inetSocketAddress2 = new InetSocketAddress("localhost", 10085);
        return Arrays.asList(inetSocketAddress, inetSocketAddress2);
    }

}
