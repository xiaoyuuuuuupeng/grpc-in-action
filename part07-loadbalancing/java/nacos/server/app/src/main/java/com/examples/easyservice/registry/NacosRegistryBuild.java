package com.examples.easyservice.registry;

import com.alibaba.nacos.api.NacosFactory;
import com.alibaba.nacos.api.exception.NacosException;
import com.alibaba.nacos.api.naming.NamingService;
import com.alibaba.nacos.api.naming.pojo.Instance;
import com.examples.easyservice.utils.NetUtils;
import io.grpc.BindableService;
import io.grpc.ServerServiceDefinition;
import org.apache.commons.lang3.StringUtils;

import java.net.URI;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;
import java.util.logging.Logger;

public class NacosRegistryBuild {
    private final Logger logger = Logger.getLogger("com.examples.easyservice.registry.NacosRegistryBuild");

    private URI uri;

    private Integer servicePort;

    private NamingService namingService = null;

    private BindableService[] bindableServices;


    public NacosRegistryBuild(URI uri,Integer servicePort,BindableService[] bindableServices) {
        this.uri = uri;
        this.bindableServices = bindableServices;
        this.servicePort = servicePort;
    }

    public void init(){
            Properties properties = buildNacosProperties(uri);
            try {
                this.namingService = NacosFactory.createNamingService(properties);
            } catch (NacosException var8) {
                this.logger.info("init grpc server failed, error: "+ var8.getErrMsg());
            }

            BindableService[] var4 = bindableServices;
            int var5 = bindableServices.length;

            for(int var6 = 0; var6 < var5; ++var6) {
                BindableService bindableService = var4[var6];
                this.registerService(bindableService);
            }

    }

    private void registerService(BindableService bindableService) {
        ServerServiceDefinition serverServiceDefinition = bindableService.bindService();
        try {
            namingService.registerInstance(serverServiceDefinition.getServiceDescriptor().getName(), createInstance());
        } catch (Exception e) {
            logger.info("Register grpc service error "+ e.getMessage());
        }
    }

    public static Properties buildNacosProperties(URI uri) {
        Properties properties = new Properties();
        StringBuilder serverAddrBuilder = (new StringBuilder(uri.getHost())).append(":").append(uri.getPort());
        String serverAddr = serverAddrBuilder.toString();
        properties.put("serverAddr", serverAddr);
        return properties;
    }

    private Instance createInstance() {
        Instance instance = new Instance();
        instance.setIp(NetUtils.getLocalHost());
        instance.setPort(servicePort);
        return instance;
    }

}
