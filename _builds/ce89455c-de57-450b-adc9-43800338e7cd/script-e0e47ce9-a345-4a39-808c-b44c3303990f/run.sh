set -e
go get github.com/xchapter7x/versioning
export NEXT_VERSION=`versioning bump_patch`
echo "next version should be: ${NEXT_VERSION}"
