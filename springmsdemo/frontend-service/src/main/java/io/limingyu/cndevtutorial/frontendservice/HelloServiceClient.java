package io.limingyu.cndevtutorial.frontendservice;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;

@FeignClient(name = "hello-service", url = "${HELLO_SERVICE_URL:}")
public interface HelloServiceClient {
    @GetMapping("/")
    String getHello();
}
