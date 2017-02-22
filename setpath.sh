#!/bin/sh

# To make current directory to GOPATH
export GOPATH=$PWD

# For pcre library header
export CGO_CFLAGS=-I/usr/local/include

# Dictionary resource
export KRGO_DIC_RSRC=$PWD/src/resources/

echo $GOPATH
echo $CGO_CFLAGS
echo $KRGO_DIC_RSRC
