# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base
FROM nginx:latest
COPY ./build /usr/share/nginx/html
ADD default.conf /etc/nginx/conf.d/default.conf
