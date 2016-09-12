#!/bin/sh

export DEBIAN_FRONTEND=noninteractive

ACTION=$1

PACKAGE=$2

if [ "$(id -u)" != "0" ]; then
	sudo apt-get $ACTION -y -q -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" $PACKAGE
else
	apt-get $ACTION -y -q -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" $PACKAGE
fi
