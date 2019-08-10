#!/bin/bash

# Calls the build script you created to rebuild the API server Linux executable and API docker container image
./build.sh

# Pushes your API server Docker container image to Docker Hub
docker push demitu/gateway

ssh ec2-user@ec2-18-224-239-175.us-east-2.compute.amazonaws.com 'bash -s' < run.sh