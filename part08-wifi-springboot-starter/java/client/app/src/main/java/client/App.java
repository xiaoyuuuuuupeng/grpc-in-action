
package client;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.Resource;
import java.util.UUID;

@SpringBootApplication
@RestController
public class App {
    @Resource
    private Client client;

    public static void main(String[] args) {
        SpringApplication.run(App.class,args);
    }

    @RequestMapping("/test")
    public String test(){
        String s = UUID.randomUUID().toString();
        String casLogin = client.casLogin(s);
        System.out.println(casLogin);
        return casLogin;
    }
}
