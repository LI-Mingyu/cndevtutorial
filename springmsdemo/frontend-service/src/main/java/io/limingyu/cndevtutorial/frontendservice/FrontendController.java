package io.limingyu.cndevtutorial.frontendservice;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class FrontendController {
    @Autowired
    private HelloServiceClient helloServiceClient;

    @GetMapping("/")
    public String greeting() {
        return helloServiceClient.getHello();
    }
}