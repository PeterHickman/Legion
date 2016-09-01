#!/bin/sh

export DEBIAN_FRONTEND=noninteractive

apt-get update -q
apt-get dist-upgrade -q -y
