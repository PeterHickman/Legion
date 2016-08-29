#!/bin/bash

TABLES_ROOT='/etc/iptables'

BLACKLIST_FILE=$TABLES_ROOT/blacklist.txt

SERVICES_FILE=$TABLES_ROOT/services.txt

##
# These need to be true and there is no corrective action that we
# can take if they are not
##
check()
{
  echo "Checking prerequisites"

  if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root" 1>&2
    exit 1
  else
    echo "- We are running as root"
  fi

  FOUND=`which iptables`
  if [ -z "$FOUND" ]; then
    echo "iptables does not appear to be installed" 1>&2
    exit 1
  else
    echo "- iptables is installed"
  fi

  FOUND=`uname -s`
  if [ "$FOUND" != "Linux" ]; then
    echo "This software assumes we are running Linux" 1>&2
    exit 1
  else
    echo "- We are running on Linux"
  fi
}

##
# Things that need to exist but we can fix them if they do not
##
needed()
{
  echo "Things that need to be set up"

  if [ ! -d "$TABLES_ROOT" ]; then
    echo "- Creating $TABLES_ROOT (have you installed iptables-persistent yet?)"
    mkdir $TABLES_ROOT
  else
    echo "- $TABLES_ROOT exists"
  fi

  # Force the permissions
  chmod a=rx,u+w $TABLES_ROOT

  if [ ! -r "$BLACKLIST_FILE" ]; then
    echo "- Creating an empty $BLACKLIST_FILE"
    touch $BLACKLIST_FILE
  else
    echo "- $BLACKLIST_FILE exists"
  fi

  # Force the permissions
  chmod a=r,u+w $BLACKLIST_FILE

  if [ ! -r "$SERVICES_FILE" ]; then
    echo "- Creating a $SERVICES_FILE file with port 22 set"
    echo 22 > $SERVICES_FILE
  else
    echo "- $SERVICES_FILE exists"
  fi

  # Force the permissions
  chmod a=r,u+w $SERVICES_FILE
}

check
needed

install -g root -o root -m a=r,u+x fw /usr/local/sbin/fw

echo "Done"
