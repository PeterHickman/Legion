#!/bin/sh

HOSTNAME=$1

hostname $HOSTNAME

echo $HOSTNAME > /etc/hostname

echo Hostname is $HOSTNAME
