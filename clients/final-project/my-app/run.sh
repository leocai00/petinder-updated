#!/bin/bash

export TLSCERT=/etc/letsencrypt/live/leo00.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/leo00.me/privkey.pem

docker rm -f petinder
docker pull leocai001/petinder
docker run -d \
--name petinder \
-p 443:443 \
-p 80:80 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
leocai001/petinder
exit