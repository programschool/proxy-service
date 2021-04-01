FROM ubuntu:20.04

RUN rm -f /etc/localtime
RUN ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY entry-proxy /home
COPY config.json /home
WORKDIR /home

# CMD ["java", "-jar", "backend-0.0.1.jar", "--spring.profiles.active=online"]
