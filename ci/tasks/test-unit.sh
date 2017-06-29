#!/bin/bash -eux

mkdir -p $GOPATH/src/github.com/blabbertabber/
cp -Rp src/ /go/src/github.com/blabbertabber/speechbroker
cd $GOPATH/src/github.com/blabbertabber/speechbroker
go get ./...
ginkgo -v -r .
