#!/bin/sh

##
# Fix things incase the install broke
##
dpkg --configure -a

apt-get update
apt-get dist-upgrade -y
