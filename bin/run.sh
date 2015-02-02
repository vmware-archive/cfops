#!/bin/bash

set -e

~/go/src/github.com/pivotalservices/cfops/out/cfops b --host 10.9.8.30 -u admin -p admin --tp tempest -d /tmp/cfops
