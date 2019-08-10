#!/bin/bash

# Builds the Docker container
docker build -t demitu/pet .

# Pushes your API server Docker container image to Docker Hub
docker push demitu/pet