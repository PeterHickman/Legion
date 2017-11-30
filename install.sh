#!/bin/sh

BINARY='/usr/local/bin/'
MAN='/usr/local/share/man/man1/'

echo "Installing legion to $BINARY"

install -v legion $BINARY

echo "Installing man page to $MAN"

install -v legion.1 $MAN
