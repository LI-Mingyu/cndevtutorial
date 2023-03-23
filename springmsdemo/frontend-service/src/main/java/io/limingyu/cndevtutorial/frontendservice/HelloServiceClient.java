package io.limingyu.cndevtutorial.frontendservice;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;

@FeignClient(name = "hello-service", url = "${HELLO_SERVICE_URL:}")
public interface HelloServiceClient {
    @GetMapping("/")
    String getHello();
}
// 如果在k8s中，上面注解中的url课通过configmap注入环境变量，如"http://hello-service"
// 若在非k8s环境且无环境变量则用默认值（冒号后的）为空，将从Eureka中获取服务地址
// url 也可以用占位符，写在properties文件中，如果部署在k8s中，用configmap覆盖properties文件
// 同样 properties占位符的值可以为空，若非K8s环境，则OpenFeign会从Eureka中获取服务地址