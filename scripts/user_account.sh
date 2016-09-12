#!/bin/sh

USER=$1

FOUND=`grep ^$USER: /etc/passwd`

if [ ! -z "$FOUND" ]; then
  echo "The account for user $USER already exists"
else
  echo "Creating $USER user"

  adduser --home /home/$USER --gecos "" --disabled-password $USER
fi

echo "Checking group membership"
adduser $USER sshlogin

KEY_FILE="/home/$USER/.ssh/id_rsa"

if [ -r "$KEY_FILE" ]; then
  echo "SSH keygen already run"
else
  echo "Setting up SSH keygen"
  su -l $USER -c "ssh-keygen -q -t rsa -N '' -f $KEY_FILE"
fi

AUTH_FILE="/home/$USER/.ssh/authorized_keys"

cat authorized_keys > $AUTH_FILE
chown $USER:$USER $AUTH_FILE
chmod a=r,u+w $AUTH_FILE
rm authorized_keys
