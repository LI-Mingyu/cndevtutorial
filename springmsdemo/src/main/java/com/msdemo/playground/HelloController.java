package com.msdemo.playground;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.PostMapping;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;


@RestController
public class HelloController {
    private static final Logger logger = LogManager.getLogger("Handling HTTP request");

    @GetMapping("/")
    public String hello() {
        logger.info("GET");
        return "Hello Spring!\n";
    } 

    @PostMapping("/")
    public String aPost() {
        logger.info("POST");
        return "Hello by POST.\n";
    }
}