FROM maven:3.8.2-jdk-11 AS build-env
ADD src /app/src
ADD pom.xml /app/
WORKDIR /app
RUN mvn package 
EXPOSE 8080
CMD ["/app/target/msdemo-1.0.jar"] 
