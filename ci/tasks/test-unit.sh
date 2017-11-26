#!/bin/bash -eux

mkdir -p $GOPATH/src/github.com/blabbertabber/
cp -Rp speech_broker/ /go/src/github.com/blabbertabber/speechbroker
cd $GOPATH/src/github.com/blabbertabber/speechbroker
go get ./...
ginkgo -v -r .
