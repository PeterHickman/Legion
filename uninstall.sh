#!/bin/sh

BINARY='/usr/local/bin/'
MAN='/usr/local/share/man/man1/'

echo "Removing legion from $BINARY"

rm -f $BINARY/legion

echo "Removing man page from $MAN"

rm -f $MAN/legion.1
