#!/bin/sh

ACTION=$1

PACKAGE=$2

apt-get $ACTION $PACKAGE -y -q
