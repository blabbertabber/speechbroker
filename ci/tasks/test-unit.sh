#!/bin/bash -eux

mkdir -p $GOPATH/src/github.com/blabbertabber/
cp -Rp src/ /go/src/github.com/blabbertabber/DiarizerServer
cd $GOPATH/src/github.com/blabbertabber/DiarizerServer
ginkgo -v -r .
