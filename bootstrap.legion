# This file bootstraps a box installing the basics to
# run and secure the system. Some of the principals are
# general but the specifics are mine

# This first step means that all the software that we require
# to complete this script is available and up to date before we
# get down to the hard work

run scripts/updates.sh

# We install ruby because some of the scripts that Legion uses are written in ruby
run scripts/package.sh install ruby

# All of these are just tools that I like to work with
run scripts/package.sh install runit
run scripts/package.sh install wget
run scripts/package.sh install lynx
run scripts/package.sh install vim
run scripts/package.sh install htop
run scripts/package.sh install curl
run scripts/package.sh install traceroute
run scripts/package.sh install lsof
run scripts/package.sh install nscd

run scripts/bootstrap/timezone.sh {timezone}
run scripts/bootstrap/hostname.sh {host}
run scripts/bootstrap/etc_hosts.rb

run scripts/bootstrap/sshlogin_group.sh
copy scripts/bootstrap/files/authorized_keys /root/authorized_keys
run scripts/bootstrap/admin_user.sh {admin}

copy scripts/bootstrap/files/sudoers /etc/sudoers
ex chmod ug=r /etc/sudoers

copy scripts/bootstrap/files/sshd_config /etc/ssh/sshd_config
ex chmod a=r,u+w /etc/ssh/sshd_config

call firewall.l

ex reboot
