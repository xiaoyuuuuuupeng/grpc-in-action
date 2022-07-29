package com.example.casservice.resolve;

import io.grpc.NameResolver;
import io.grpc.NameResolverProvider;
import java.net.URI;


public class CustomNameResolverProvider extends NameResolverProvider {

    public CustomNameResolverProvider() {

    }
    @Override
    public NameResolver newNameResolver(URI targetUri, NameResolver.Args args) {
        return new CustomNameResolver(targetUri.getPath());
    }

    @Override
    protected boolean isAvailable() {
        return true;
    }

    @Override
    protected int priority() {
        return 10;
    }

    @Override
    public String getDefaultScheme() {
        return "example";
    }
}
