#!/bin/sh

if [ "$(id -u)" != "0" ]; then
  SUDO_PREFIX='sudo'
else
  SUDO_PREFIX=''
fi

export DEBIAN_FRONTEND=noninteractive

$SUDO_PREFIX apt-get update -q
$SUDO_PREFIX apt-get dist-upgrade -q -y
