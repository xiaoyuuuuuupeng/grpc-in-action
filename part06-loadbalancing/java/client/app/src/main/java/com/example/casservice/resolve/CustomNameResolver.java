
package com.example.casservice.resolve;

import com.google.common.collect.Lists;
import io.grpc.EquivalentAddressGroup;
import io.grpc.NameResolver;

import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

import io.grpc.Attributes;


public class CustomNameResolver extends NameResolver {

    private final String path;

    private Listener2 listener;


    public CustomNameResolver(String path) {
        this.path = path;
    }

    @Override
    public String getServiceAuthority() {
        return this.path;
    }

    @Override
    public void shutdown() {
    }

    @Override
    public void start(Listener2 listener) {
        this.listener = listener;
        this.resolve();
    }



    private void resolve() {
        List<SocketAddress> socketAddressList = getAddressList(path).stream()
                .map(this::toSocketAddress)
                .collect(Collectors.toList());

        List<EquivalentAddressGroup> equivalentAddressGroups = Lists.newArrayList();
        for (SocketAddress socketAddress : socketAddressList) {
            EquivalentAddressGroup equivalentAddressGroup = new EquivalentAddressGroup(socketAddress);
            equivalentAddressGroups.add(equivalentAddressGroup);
        }
        ResolutionResult resolutionResult = ResolutionResult.newBuilder()
                .setAddresses(equivalentAddressGroups)
                .setAttributes(Attributes.EMPTY)
                .build();

        this.listener.onResult(resolutionResult);

    }
    private SocketAddress toSocketAddress(InetSocketAddress address) {
        return new InetSocketAddress(address.getHostName(), address.getPort());
    }


    private List<InetSocketAddress> getAddressList(String path) {
        System.out.println("path:"+path);
        InetSocketAddress inetSocketAddress = new InetSocketAddress("localhost", 10084);
        InetSocketAddress inetSocketAddress2 = new InetSocketAddress("localhost", 10085);
        return Arrays.asList(inetSocketAddress, inetSocketAddress2);
    }

}
