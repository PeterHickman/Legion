#!/bin/sh

USER=$1

FOUND=`grep ^$USER: /etc/passwd`

if [ "$(id -u)" != "0" ]; then
  SUDO_PREFIX='sudo'
else
  SUDO_PREFIX=''
fi

if [ ! -z "$FOUND" ]; then
  echo "The account for user $USER already exists"
else
  echo "Creating $USER user"

  $SUDO_PREFIX adduser --home /home/$USER --gecos "" --disabled-password $USER
fi

echo "Checking group membership"
$SUDO_PREFIX adduser $USER sshlogin

KEY_FILE="/home/$USER/.ssh/id_rsa"

FOUND=`$SUDO_PREFIX [ -r "$KEY_FILE" ] && echo 'yes'`

if [ "$FOUND" = "yes" ]; then
  echo "SSH keygen already run"
else
  echo "Setting up SSH keygen"
  $SUDO_PREFIX su -l $USER -c "ssh-keygen -q -t rsa -N '' -f $KEY_FILE"
fi

AUTH_FILE="/home/$USER/.ssh/authorized_keys"

$SUDO_PREFIX su -l $USER -c "cat /tmp/authorized_keys > $AUTH_FILE"
$SUDO_PREFIX chown $USER:$USER $AUTH_FILE
$SUDO_PREFIX chmod a=r,u+w $AUTH_FILE
rm /tmp/authorized_keys
