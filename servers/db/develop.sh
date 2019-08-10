#!/bin/bash

# Builds the Docker container
docker build -t demitu/db .

# Pushes your API server Docker container image to Docker Hub
docker push demitu/db