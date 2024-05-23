#!/bin/sh

BINARY='/usr/local/bin'

echo "Building legion"
go build legion.go

echo "Installing legion to $BINARY"
install -v legion $BINARY

echo "Removing the build"
rm legion
