FROM nginx:1.18-alpine-perl
COPY index.html /usr/share/nginx/html/
RUN echo $(date -R) >> /usr/share/nginx/html/index.html
