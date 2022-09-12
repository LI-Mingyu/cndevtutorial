package io.limingyu.cndevtutorial.frontendservice;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.cloud.client.ServiceInstance;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import java.util.List;

@RestController
public class FrontendController {
    private static final Logger logger = LogManager.getLogger("Handling HTTP request");

    @Autowired
    private DiscoveryClient discoveryClient;

    @GetMapping("/")
    public String greeting() {
        List<ServiceInstance> list = discoveryClient.getInstances("hello-servcie");
        if (list != null && list.size() > 0) {
            RestTemplate restTemplate = new RestTemplate();
            ResponseEntity<String> responseEntity = restTemplate.getForEntity(list.get(0).getUri(), String.class);
            return responseEntity.getBody();
        }
        return "";
    }
}