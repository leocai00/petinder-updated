#!/bin/bash
./build.sh

docker push leocai001/petinder

ssh -i ~/.ssh/id_rsa ec2-user@ec2-18-220-69-125.us-east-2.compute.amazonaws.com 'bash -s' < run.sh