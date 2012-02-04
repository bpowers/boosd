#!/bin/sh

pushd parser
make
popd

go install
