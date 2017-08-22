#!/bin/bash -ue

# Install git
yum install -y git

#Install GO Lang
cd /tmp/
curl https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz | tar -C /usr/local -xzf -
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile.d/path.sh

source /etc/profile.d/path.sh
export GOPATH="/var/tmp/"

#Get the source code
go get github.com/pivotalservices/cfops || true

#Install Glide
export GOBIN="/usr/local/go/bin"
curl https://glide.sh/get | sh

#Pull in glide managed dependencies:
cd $GOPATH/src/github.com/pivotalservices/cfops
glide install

#Build the project:
cd cmd/cfops/
go build

ls -la /var/tmp/src/github.com/pivotalservices/cfops/cmd/cfops/cfops
