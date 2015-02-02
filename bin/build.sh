#!/bin/bash

set -e

echo -e "\nGenerating Binary..."

ROOT_DIR=$(cd $(dirname $(dirname $0)) && pwd)

echo -e "\nROOT_DIR set to ${ROOT_DIR}"

CLI_GOPATH=$ROOT_DIR/tmp/cli_gopath
rm -rf $CLI_GOPATH
mkdir -p $CLI_GOPATH/src/github.com/pivotalservices/
ln -s $ROOT_DIR $CLI_GOPATH/src/github.com/pivotalservices/cfops

GODEP_GOPATH=$ROOT_DIR/Godeps/_workspace

GOPATH=$CLI_GOPATH:$GODEP_GOPATH:$GOPATH go build -o $ROOT_DIR/out/cfops cmd/cfops/main.go
rm -rf $CLI_GOPATH

CFOPS_HOME=~/.cfops
CFOPS_CONFIG=${CFOPS_HOME}/config.json
if [[ ! -f $CFOPS_CONFIG ]]; then
	mkdir -p ${CFOPS_HOME}
	echo -e "Copying configuration to ${CFOPS_HOME}"
	cp $ROOT_DIR/config/assets/config.json ${CFOPS_HOME}
fi
