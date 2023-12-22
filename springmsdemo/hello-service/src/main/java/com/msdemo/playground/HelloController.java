package com.msdemo.playground;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.PostMapping;
import io.github.resilience4j.ratelimiter.annotation.RateLimiter;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;


@RestController
public class HelloController {
    private static final Logger logger = LogManager.getLogger("Handling HTTP request");

    @GetMapping("/")
    @RateLimiter(name = "helloService", fallbackMethod = "rateLimitFallback")
    public String hello() {
        logger.info("GET");
        return "Hello Spring!\n";
    } 

    @GetMapping("/{id}") 
    public String hello(@PathVariable("id") String id) {
        logger.info("GET: " + id);
        return "Hello " + id + "!\n";
    }

    @PostMapping("/")
    public String aPost() {
        logger.info("POST");
        return "Hello by POST.\n";
    }

    public String rateLimitFallback(Throwable t) {
        logger.info("Rate limit fallback");
        return "收银台正在排队，请稍后再试.\n";
    }
}