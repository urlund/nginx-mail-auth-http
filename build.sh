#!/bin/bash
VERSION=`cat VERSION`
GOOS=linux
GOARCH=amd64
BIN_NAME="nginx-mail-auth-http"
BIN_PATH="bin/$VERSION/$BIN_NAME-$GOOS-$GOARCH"
GIT_REPO="github.com/urlund/$BIN_NAME"

rm -rf $BIN_NAME

sed -i '' -e "s/\"version: .*\"/\"version: $VERSION\"/g" init.go
sed -i '' -e "s/download\/.*\/nginx/download\/$VERSION\/nginx/g" Dockerfile

# check if docker is installed
if [ $(which docker > /dev/null 2>&1; echo $?) -ne 0 ]; then
    echo "you must have docker installed to use this script"
    exit 1;
fi

# check if $GOPATH isset
if [ -z $GOPATH ]; then
    echo "you must set \$GOPATH to use this script"
    exit 1;
fi

# run docker with build cmd
docker run --rm -it -v "$GOPATH":/work -e "GOPATH=/work" -w /work/src/$GIT_REPO -e GOOS=$GOOS -e GOARCH=$GOARCH golang:latest go build -o $BIN_PATH
