#!/bin/sh

ZONE=$1

echo $ZONE > /etc/timezone
rm /etc/localtime
ln -s /usr/share/zoneinfo/$ZONE /etc/localtime

echo Timezone is now $ZONE
