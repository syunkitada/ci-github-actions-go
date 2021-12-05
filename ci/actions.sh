#!/bin/bash -ex

COMMAND="${@:-help}"

NAME=${PWD##*/}
WD=${PWD}
MODULE=`grep '^module ' go.mod | awk '{print $2}'`
GOPATH=/opt/ci/go-1.17.4
MODULE_PATH=/opt/ci/go-1.17.4/${MODULE}

function setup_go() {
	mkdir -p /opt/ci
	test -e /opt/ci/go || \
        (wget https://go.dev/dl/go1.17.4.linux-amd64.tar.gz && tar xvf go1.17.4.linux-amd64.tar.gz -C /opt/ci/ && rm -f go1.17.4.linux-amd64.tar.gz)

    mkdir -p /opt/ci/gcov2lcov
	test -e /opt/ci/gcov2lcov/bin || \
        (wget https://github.com/jandelgado/gcov2lcov/releases/latest/download/gcov2lcov-linux-amd64.tar.gz && \
         tar xvf gcov2lcov-linux-amd64.tar.gz -C /opt/ci/gcov2lcov && rm -f gcov2lcov-linux-amd64.tar.gz)

    test -e $GOPATH/src/$MODULE || \
        (mkdir -p $GOPATH/src/$MODULE && cp -r $WD/. $GOPATH/src/$MODULE)

    cat << EOS > /opt/ci/envgorc
PATH=$PATH:/opt/ci/go:/opt/ci/gcov2lcov
GOPATH=$GOPATH
MODULE_PATH=$MODULE_PATH
EOS
}

function test_go() {
	. /opt/ci/envgorc
    cd ${MODULE_PATH}
    go test --coverpkg ./pkg/... -coverprofile=.coverage.out ./pkg/...
	gcov2lcov-linux-amd64 -infile .coverage.out -outfile /tmp/coverage.lcov
}

$COMMAND
