#!/bin/sh

export DEBIAN_FRONTEND=noninteractive

ACTION=$1

PACKAGE=$2

apt-get $ACTION -y -q -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" $PACKAGE
