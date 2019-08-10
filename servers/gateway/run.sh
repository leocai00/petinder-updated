docker network disconnect slack gateway
docker network disconnect slack db
docker network disconnect slack redis
docker network disconnect slack messaging
docker network disconnect slack mongo
docker network disconnect slack rabbit
docker network disconnect slack pet
docker network disconnect slack mongopet

# Stop and remove the existing container instance
docker rm -f gateway
docker rm -f db
docker rm -f redis
docker rm -f messaging
docker rm -f mongo
docker rm -f mongopet
docker rm -f rabbit
docker rm -f pet

docker network rm slack

# Pull the updated container image from DockerHub
docker pull demitu/gateway
docker pull demitu/db
docker pull demitu/messaging
docker pull demitu/pet
docker pull mongo

docker network create slack

# Export variables
export TLSCERT=/etc/letsencrypt/live/api.demitu.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.demitu.me/privkey.pem

# Re-run a newly-updated instance
docker run -d --name redis --network slack --restart on-failure redis
docker run -d --name mongo --network slack --restart on-failure mongo
docker run -d --name mongopet --network slack --restart on-failure mongo
sleep 30
docker run -d --name rabbit --network slack --restart on-failure rabbitmq:3-management
sleep 30
docker run -d --name messaging --network slack --restart on-failure -e RABBITADDR=rabbit:5672 -e MONGO=mongo:27017 demitu/messaging
docker run -d --name pet --network slack --restart on-failure -e RABBITADDR=rabbit:5672 -e MONGO=mongopet:27017 demitu/pet
docker run -d --name db --network slack --restart on-failure -e MYSQL_ROOT_PASSWORD="test" -e MYSQL_DATABASE=db demitu/db
docker run -d --name gateway --network slack --restart on-failure -p 443:443 -e SESSIONKEY=$SESSIONKEY -e REDISADDR=redis:6379 -e MESSAGEADDR=messaging:80 -e PETADDR=pet:80 -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e MYSQL_ROOT_PASSWORD="test" -e DBADDR=db:3306 -e RABBITADDR=rabbit:5672 -v /etc/letsencrypt:/etc/letsencrypt:ro demitu/gateway