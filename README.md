# Legion - configuration and deployment software

Legion allows the automated deployment and configuration of machines in the vein of
Chef or Ansible but it allows me to use my existing knowledge for configuring machines
rather than having to learn a whole new framework

It was motivated by my hosting provider going out of business and me having to deploy
to a new server. Thats when I realised that I would have to do it all by hand :(

Why is it called "Legion"? Well I have created several such system before of much
complexity that ultimately made it harder to deploy despite their sweet architecture. 
I've written so many of them they are legion :)

This time I have gone for the simplest solution

## Prerequisites

### ssh is installed and running on the target machine

Generally speaking this is a given but just to be sure I'll state it here

### root can log in via ssh with a password

This is not always the case. Root usually does not have a password and the `/etc/ssh/sshd_config`
file needs to have the following directive

    PermitRootLogin yes

This is a vulnerability but after the configuration it will be closed

Also create a password for root if you don't have one

### ruby is installed

In most cases this is true but there are some systems that do not have Ruby installed
(the Raspberry Pi is an example of such a system). There is no great reliance on Ruby
-- its just something that I use a lot -- so it could be swapped out for Python in the 
future as it is pretty much guaranteed to be available

### apt-get is the package manager

This is just a statement of fact rather than a requirement. It wouldn't take much to
have the system package manager agnostic -- using `yum` for example -- I just haven't
had a reason to make the change yet

## Warning!!!

This is a work in progress and things will change

## Using Legion

Lets use the `bootstrap.txt` script as a worked example. The pupoose of the script is to take a newly provisioned machine and configure it to a known baseline. Before we start to run the script the target machine will be:

0. Running ssh
1. Allow `root` to log in using a password
2. `root` will have a password

After the script has been run the target machine will have:

0. Only allow access via ssh by key exchange
1. The only, initial, access will be to a named admin account
2. The firewall is up and running allowing access only to port 22
3. The hostname will be set
4. The timezone will be set

Ok. Given that the `bootstrap.txt` script is a general purpose script we need to put the target specific configuration information in it's own file -- called `server.txt` in this example

```
# Describes the server

set host pi1.local
set timezone Europe/London
set admin fred
```

This file just defines three variables that descibe the target configuration information, the hostname, timezone and the name of the admin account. With this the bootstrap script can simply refer to the variable names.

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

The first step is to make sure that the existing packages are up to date and then install the various packages that we will need. Of all the packages `ruby` is the only essential one as it is required later in the bootstrap process.

Commands like `run scripts/updates.sh` tell Legion to upload the file `updates.sh` from the `scripts` directory that came with Legion to the target machine and run it there. `updates.sh` itself is just a bash script.

```
run scripts/bootstrap/sshlogin_group.sh

run scripts/bootstrap/timezone.sh {timezone}

run scripts/bootstrap/hostname.sh {host}

run scripts/bootstrap/etc_hosts.rb
```

More scipts are uploaded and run to create the `sshlogin` group (members of this group are the only once who can log in via ssh), set the timezone and hostname and finally make sure that the `/etc/hosts` file correctly reflects the hostname.

Of note here is the `{timezone}` argument to the timezone configuration. In the `server.txt` script earlier we set `timezone` to `Europe/London`. When Legion encounters the line in the `bootstrap.txt` script it will replace `{timezone}` with `Europe/London` before running the scipt to the target machine. This allows us to have a machine specific machine configuration without having to duplicate the whole of the bootstrap script with just a few minor changes.

```
copy scripts/bootstrap/files/authorized_keys2 /root/authorized_keys2
run scripts/bootstrap/admin_user.sh {admin}
```

Legion copies the `authorized_keys2` to the target machine and then creates the admin user

```
copy scripts/bootstrap/files/sudoers /etc/sudoers
ex chmod ug=r /etc/sudoers

copy scripts/bootstrap/files/sshd_config /etc/ssh/sshd_config
ex chmod a=r,u+w /etc/ssh/sshd_config

@ firewall.txt

ex reboot
```

Finishing off we copy up our customised sudoers file and make sure the permissions are set correctly. The `ex` command executes the rest of the line directly on the target machine.

After doing the same for the `sshd_config` file we configure the firewall. The firewall configuration in `firewall.txt` is just another Legion script but rather than fill this script with it's verbage we put it in it's own file and then call that instead. `@` can be used to call other Legion scripts to make things a little modular

## Testing the scripts

Because a mistake can result in the target machine becoming unusable you will need to test the scripts before deployment. I personally use these methods:

0. Set up a Raspberry Pi with either Debian or CentOS
1. User Oracle VirtualBox to run the target OS
2. Set up a vm on Rackspace

Basically test everything before deploying, then test it some more

## Running Legion
