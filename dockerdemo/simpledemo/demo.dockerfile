FROM centos:7
RUN yum update -y; yum install -y epel-release; yum install -y nginx
COPY index.html /usr/share/nginx/html/
RUN echo $(date -R) >> /usr/share/nginx/html/index.html
EXPOSE 80
CMD "nginx" "-g" "daemon off;"
