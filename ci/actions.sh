#!/bin/bash -ex

COMMAND="${@:-help}"
NAME=${PWD##*/}
WD=${PWD}

function setup_go() {
    GOPATH=/opt/ci/go-1.17.4
    MODULE=`grep '^module ' go.mod | awk '{print $2}'`
    MODULE_PATH=/opt/ci/go-1.17.4/src/${MODULE}

	mkdir -p /opt/ci
	test -e /opt/ci/go || \
        (wget -q https://go.dev/dl/go1.17.4.linux-amd64.tar.gz && tar xf go1.17.4.linux-amd64.tar.gz -C /opt/ci/ && rm -f go1.17.4.linux-amd64.tar.gz)

    test -e $MODULE_PATH || \
        (mkdir -p $MODULE_PATH && cp -r $WD/. $MODULE_PATH)

    cat << EOS > /opt/ci/envgorc
export PATH=/opt/ci/go/bin:${GOPATH}/bin:$PATH
export GOROOT=/opt/ci/go
export GOPATH=${GOPATH}
export GO111MODULE=off
export MODULE_PATH=$MODULE_PATH
EOS

	. /opt/ci/envgorc
    go version
    cd $GOROOT
    GO111MODULE=on

    # for coveralls
    go install github.com/jandelgado/gcov2lcov@latest

    # for codecov
    go install github.com/AlekSi/gocov-xml@latest
}

function test_go() {
	. /opt/ci/envgorc
    cd ${MODULE_PATH}
    go test --coverpkg ./pkg/... -coverprofile=.coverage.out ./pkg/...

    # for coveralls
    gcov2lcov -infile .coverage.out -outfile /tmp/coverage.lcov

    # for codecov
    gocov convert .coverage.out | gocov-xml > /tmp/coverage.xml
}

$COMMAND
