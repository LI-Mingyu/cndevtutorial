FROM maven:3.8.2-jdk-11 AS build-env
ADD src /app/src
ADD pom.xml /app/
WORKDIR /app
RUN mvn package 

FROM gcr.io/distroless/java:11
COPY --from=build-env /app/target/frontend-service-0.0.1-SNAPSHOT.jar /app/
WORKDIR /app
EXPOSE 8080
CMD ["frontend-service-0.0.1-SNAPSHOT.jar"]
