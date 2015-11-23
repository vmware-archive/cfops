set -e
cd $WERCKER_SOURCE_DIR
go version
go get -u github.com/tools/godep
export PATH=$WERCKER_SOURCE_DIR/bin:$PATH
