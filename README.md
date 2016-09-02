# Legion - configuration and deployment software

Legion allows the automated deployment and configuration of machines in the vein of Chef or Ansible but it allows me to use my existing knowledge for configuring machines rather than having to learn a whole new framework

It was motivated by my hosting provider going out of business and me having to deploy to a new server. Thats when I realised that I would have to do it all by hand :(

Why is it called "Legion"? Well I have created several such system before of such complexity that it made it harder to deploy despite their sweet architecture. I've written so many of them they are legion :)

This time I have gone for the simplest solution

## TL;DR

Legion is a tool that allows a simple scripting language (with a mini templating system) to run over SSH to a target system. The scripting language can issue commands, copy files from the host to the target and upload and execute stored scripts. I use it to automate the configuration of servers. It probably has other uses and there are certainly a million and one other (better) tools that do the same thing

## Prerequisites

### ssh is installed and running on the target machine

Generally speaking this is a given but just to be sure I'll state it here

### root can log in via ssh with a password

This is not always the case. Root usually does not have a password and the `/etc/ssh/sshd_config` file needs to have the following directive

    PermitRootLogin yes

This is bad security practice so create a good long password. Once we have bootstrapped the system it will be closed anyway

Create a password for root if you don't have one

### apt-get is the package manager

This is just a statement of fact rather than a requirement. It wouldn't take much to have the system package manager agnostic -- using `yum` for example -- I just haven't had a reason to make the change yet

## Warning!!!

This is a work in progress and things will change

## Using Legion

Lets use the `bootstrap.legion` script as a worked example. The purpose of the script is to take a newly provisioned machine and configure it to a known baseline. Before we start to run the script the target machine will be:

0. Running ssh
1. Allow `root` to log in using a password
2. `root` has a password

After the script has been run the target machine will:

0. Only allow access via ssh by key exchange
1. The initial access will be by the named admin account
2. The firewall is up and running allowing access only to port 22
3. The hostname will be set
4. The timezone will be set

Ok. Given that the `bootstrap.legion` script is a general purpose script we need to put the target specific configuration in it's own file -- called `server.legion` in this example (which is for my Raspberry Pi)

```
# Describes the server

set host pi1.local
set timezone Europe/London
set admin fred
```

This file just defines three variables that are specific to the target's configuration. The hostname, timezone and the name of the admin account. With this the bootstrap script can simply refer to the variable names.

Now lets walk through the bootstrap script itself, less the comments.

```
run scripts/updates.sh

run scripts/package.sh install ruby
run scripts/package.sh install runit
run scripts/package.sh install wget
run scripts/package.sh install lynx
run scripts/package.sh install vim
run scripts/package.sh install htop
run scripts/package.sh install curl
run scripts/package.sh install mutt
run scripts/package.sh install traceroute
run scripts/package.sh install lsof
run scripts/package.sh install pstree
run scripts/package.sh install nscd
```

The first step is to make sure that the existing packages are up to date and then install the packages that we will need. Of all the packages `ruby` is the only essential one as it is required later in the bootstrap process.

Commands like `run scripts/updates.sh` tell Legion to upload the file `updates.sh` from the `scripts` directory that came with Legion to the target machine and run it there. `updates.sh` is an ordinary bash script.

```
run scripts/bootstrap/sshlogin_group.sh

run scripts/bootstrap/timezone.sh {timezone}

run scripts/bootstrap/hostname.sh {host}

run scripts/bootstrap/etc_hosts.rb
```

More scripts are uploaded and run to create the `sshlogin` group (members of this group are the only ones who can log in via ssh), set the timezone and hostname and finally make sure that the `/etc/hosts` file correctly reflects the hostname.

Of note here is the `{timezone}` argument to the timezone configuration. In the `server.legion` script earlier we set `timezone` to `Europe/London`. When Legion encounters the line in the `bootstrap.legion` script it will replace `{timezone}` with `Europe/London` before running the script on the target machine. This allows us to have a machine specific configuration without having to duplicate the whole of the bootstrap script with just a few minor changes for each server we want to deploy.

```
copy scripts/bootstrap/files/authorized_keys /root/authorized_keys
run scripts/bootstrap/admin_user.sh {admin}
```

Legion copies the `authorized_keys` to the target machine and then runs the script to create the admin user

```
copy scripts/bootstrap/files/sudoers /etc/sudoers
ex chmod ug=r /etc/sudoers

copy scripts/bootstrap/files/sshd_config /etc/ssh/sshd_config
ex chmod a=r,u+w /etc/ssh/sshd_config

call firewall.legion

ex reboot
```

Finishing off we copy up our customised sudoers file and make sure the permissions are set correctly. The `ex` command executes the rest of the line directly on the target machine.

After doing the same for `sshd_config` we configure the firewall. The firewall configuration in `firewall.legion` is just another Legion script but rather than fill this script with it's verbage we put it in it's own file and then call that instead. `call` can be used to call other Legion scripts to make things a little more modular

## Testing the scripts

Because a mistake can result in the target machine becoming unusable you will need to test the scripts before deployment. I personally use these methods:

0. Set up a Raspberry Pi with either Debian or CentOS
1. Use Oracle VirtualBox to run the target OS
2. Set up a vm on Rackspace

Basically test everything before deploying, then test it some more

## Running Legion

From the command line we call it thus

    $ legion --host 1.2.3.4 --port 22 --username root --password secret server.legion bootstrap.legion

All four parameters must be given even if there are sane defaults (for `port` for example) or they might not be required (`password` is not needed if you are using ssh key exchange to the admin account -- just give it some dummy value)

After the parameters are the script files that will be run in the order given. There should be at least one

## Why the long ass file extension?

All the shorter ones were already taken. Anything shorter than the full name would have been needlessly cryptic so I went the whole hog. This way Github wont think that half this project is written in Lisp
