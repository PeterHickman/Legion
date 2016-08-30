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
