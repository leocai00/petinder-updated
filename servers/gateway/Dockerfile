# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server

# Use `alpine` as the base image
FROM alpine

# Use the `RUN` command to install the `ca-certificates` package in your new container image
RUN apk add --no-cache ca-certificates

# Copy the executable that you will be running into the container
COPY gateway /gateway

# Expose port
EXPOSE 443

# Tell docker to use the executable as the main entrypoint when the
# container starts. That is, what should the container do when it starts?
ENTRYPOINT [ "/gateway" ]