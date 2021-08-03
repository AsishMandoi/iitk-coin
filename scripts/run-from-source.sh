#!/bin/bash

temp=${PWD}

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
read -p "Enter new emailid and password? (If you have entered these once, you need not provide them again) [Y/N] " bool

if ( [ "$bool" == "y" ] || [ "$bool" == "Y" ] ); then
  # Ask the user for emailid and password
  echo "Please enter the emailid and password for the sender's account from which all OTPs are to be sent."
  read -p 'Email-ID: ' emailid
  read -sp 'Password: ' password
  echo
  export EMAIL_ID=$emailid
  export PASSWORD=$password
fi

echo Running the containers \"$backend\" and \"$redis\"

######### One command to pull, build images and then run the containers
docker-compose up

cd ${temp}