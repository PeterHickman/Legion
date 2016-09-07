#!/bin/sh

GROUP=$1

FOUND=`grep ^$GROUP /etc/group`

if [ -z "$FOUND" ]
then
    addgroup $GROUP
else
    echo Group $GROUP already exists
fi
