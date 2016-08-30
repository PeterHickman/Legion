#!/bin/sh

ADMIN=$1

FOUND=`grep ^$ADMIN: /etc/passwd`

if [ ! -z "$FOUND" ]; then
  echo "The account for user $ADMIN already exists"
else
  echo "Creating $ADMIN user"

  adduser --home /home/$ADMIN --gecos "" --disabled-password $ADMIN
fi

echo "Checking group membership"
adduser $ADMIN sudo
adduser $ADMIN sshlogin

KEY_FILE="/home/$ADMIN/.ssh/id_rsa"

if [ -r "$KEY_FILE" ]; then
  echo "SSH keygen already run"
else
  echo "Setting up SSH keygen"
  su -l $ADMIN -c "ssh-keygen -q -t rsa -N '' -f $KEY_FILE"
fi

AUTH_FILE="/home/$ADMIN/.ssh/authorized_keys2"

cat authorized_keys2 > $AUTH_FILE
chown $ADMIN:$ADMIN $AUTH_FILE
chmod a=r,u+w $AUTH_FILE
rm authorized_keys2
