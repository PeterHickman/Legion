#!/bin/sh

export DEBIAN_FRONTEND=noninteractive

ACTION=$1

PACKAGE=$2

if [ "$(id -u)" != "0" ]; then
  SUDO_PREFIX='sudo'
else
  SUDO_PREFIX=''
fi

$SUDO_PREFIX apt-get $ACTION -y -q -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" $PACKAGE
