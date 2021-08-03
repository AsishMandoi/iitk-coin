#!/bin/bash

# Ask the user if they want to provide new emailid and password for the sender's account
read -p 'Enter new emailid and password? (If you have entered these once, you need not provide them again.) [Y/N] ' bool

if ( [ "$bool" == "y" ] || [ "$bool" == "Y" ] ); then
  # Ask the user for emailid and password
  echo "Please enter the emailid and password for the sender's account from which all OTPs are to be sent."
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
if [ "$(docker network inspect $network)" ]; then
  echo "Network with name \"$network\" already exists."
else
  docker network create $network
  echo "Created network \"$network\""
fi

# Function to run multiple containers on the created network
run_container(){
  if [ "$(docker ps -q -f name=$1)" ]; then
    echo "Restarting \"$1\"... If you want to stop the container, please do so by running - \"docker stop $1\"."
    docker restart $1
  else
    if [ "$(docker ps -aq -f name=$1)" ]; then
      docker rm $1
    fi
    docker run \
    -e EMAIL_ID=$EMAIL \
    -e PASSWORD=$PSWD \
    -e BACKEND_CONTAINER_NAME=$backend \
    -e REDIS_CONTAINER_NAME=$redis \
    -d --name $1 -p $2 --network $network $3

    echo "\"$1\" is running on \"$network\" in the background."
  fi
}

# Run the backend-container and the redis-container using the created network
run_container "$backend" "8080:8080" "asishmandoi/iitk-coin:latest"
run_container "$redis" "6379" "redis:6.2.5-alpine3.14"
