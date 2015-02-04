#!/bin/bash

set -e

echo -e "\nRunning cfops..."

ROOT_DIR=$(cd $(dirname $(dirname $0)) && pwd)

$ROOT_DIR/out/cfops --logLevel="debug" b --host 10.9.8.30 -u admin -p admin --tp tempest -d /tmp/cfops
