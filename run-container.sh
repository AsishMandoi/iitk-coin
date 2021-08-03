#!/bin/bash

# Ask the user if they want to provide new emailid and password for the sender's account
read -p 'Enter new emailid and password? [Y/N] ' bool

if ( [ "$bool" == "y" ] || [ "$bool" == "Y" ] ); then
  # Ask the user for emailid and password
  echo "Please enter the emailid and password of the sender's account from which all OTPs are to be sent. (This only needs to be done the first time)"
  read -p 'Email-ID: ' emailid
  read -sp 'Password: ' password
  echo
  export EMAIL=$emailid
  export PSWD=$password
fi

backend="IITK-Coin.backend.container"
redis="IITK-Coin.redis.container"

export BACKEND_CONTAINER_NAME=$backend
export REDIS_CONTAINER_NAME=$redis

# Name of the network
network="IITK-Coin.network"

# Create custom network
docker network create $network

# Run the backend-container and the redis-container using the created network
docker run -e EMAIL_ID=$EMAIL \
-e PASSWORD=$PSWD \
-e BACKEND_CONTAINER_NAME=$backend \
-e REDIS_CONTAINER_NAME=$redis \
--name $backend -p 8080:8080 -d --network $network asishmandoi/iitk-coin:latest && docker run --name $redis -p 6379 -d --network $network redis:6.2.5-alpine3.14
