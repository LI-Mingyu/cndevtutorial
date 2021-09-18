FROM centos:7
ADD https://dl.k8s.io/release/v1.21.2/bin/linux/amd64/kubectl /usr/local/bin/
RUN chmod +x /usr/local/bin/kubectl
