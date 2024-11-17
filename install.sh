#!/bin/sh

BINARY='/usr/local/bin'
APP=legion

echo "Building $APP"
go build -ldflags="-s -w" $APP.go

echo "Installing dexec to $BINARY"
install $APP $BINARY

echo "Removing the build"
rm $APP