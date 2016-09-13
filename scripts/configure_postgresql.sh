#!/bin/sh

if [ "$(id -u)" != "0" ]; then
  SUDO_PREFIX='sudo'
else
  SUDO_PREFIX=''
fi

APPLICATION_NAME=$1

##
# Find the pg_hba.conf file
##

$SUDO_PREFIX ls -1 /etc/postgresql/*/main/pg_hba.conf > x

if [ ! -s x ]; then
  echo "ERROR: The pg_hba.conf file cannot be found"
  exit 1
fi

HBA=`head -1 x`

$SUDO_PREFIX su postgres -c "psql -c 'CREATE ROLE $APPLICATION_NAME WITH NOSUPERUSER CREATEDB INHERIT LOGIN;'"

LINE="host ${APPLICATION_NAME}_production $APPLICATION_NAME 127.0.0.1/32 trust"

FOUND=`$SUDO_PREFIX grep "$LINE" $HBA`

if [ "$FOUND" != "" ]; then
  echo "pg_hba.conf already configured"
else
  echo "Configuring pg_hba.conf"
  echo "With: $LINE"
  $SUDO_PREFIX su postgres -c "echo $LINE >> $HBA"
fi

$SUDO_PREFIX /etc/init.d/postgresql reload

$SUDO_PREFIX su $APPLICATION_NAME -c "createdb ${APPLICATION_NAME}_production -E utf8"
