#!/bin/bash

######### Downloading the source code

# Make sure `GOPATH` environment variable is set
export GOPATH=$(go env GOPATH)

# Create this directory and change the current working directory to the following
mkdir -p $GOPATH/src/github.com/AsishMandoi/ && cd $GOPATH/src/github.com/AsishMandoi/

# Clone this repository
git clone https://github.com/AsishMandoi/iitk-coin.git

# Change the working directory again as follows
cd ./iitk-coin

######### Setting up required enviroment variables

# Ask the user if they want to provide new emailid and password for the sender's account
read -p 'Enter new emailid and password? [Y/N] ' bool

if ( [ "$bool" == "y" ] || [ "$bool" == "Y" ] ); then
  # Ask the user for emailid and password
  echo "Please enter the emailid and password of the sender's account from which all OTPs are to be sent. (This only needs to be done the first time)"
  read -p 'Email-ID: ' emailid
  read -sp 'Password: ' password
  echo
  export EMAIL_ID=$emailid
  export PASSWORD=$password
fi

######### One command to pull, build images and then run the containers
docker-compose up