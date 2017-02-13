# Legion - scriptable access to remote servers via SSH

Legion is a tool that allows a simple scripting language (with a mini template system) to run over SSH to a target system. The scripting language can issue commands, copy files from the host to the target and upload and execute stored scripts. I use it to automate the configuration of servers. It probably has other uses and there are certainly a million and one other (better) tools that do the same thing

Why is it called "Legion"? Well I have created several such system before of such complexity that it made it harder to deploy despite their sweet architecture. I've written so many of them they are legion :)

This time I have gone for the simplest solution

## Prerequisites

Legion requires Ruby 1.9.3 or better because it currently uses the colorize gem. I will probably just steal the parts of colorize I require at some point so that even 1.8.7 can use legion

So you need to colorize gem

    gem install colorize

## Using Legion

### Command line parameters

|Option|Required|Description|
|---|---|---|
|`--host`|**required**|The name of the host to connect to|
|`--port`|*optional*|The port to connect to the host on, defaults to `22`|
|`--username`|**required**|The username to connect to the host with|
|`--password`|**required**|The password that goes with the username, even when using ssh key exchange there needs to be a value for this, `dummy` would be fine|
|`--pretend`|*optional*|Takes no argument. Runs the script but does not actually make the `ssh` and `sftp` connections and fakes the interaction|
|`--set`|*optional*|Takes **two** arguments and acts like the `set` command within a `legion` script|
|`--log`|*optional*| Sets the log file name to whatever is given as the parameter. By default the log file will be called `legion.YYYYMMDD-HHMM.log`|
|`--config`|*optional*|An optional config file to save having really long command lines|

Parameters can be supplied as either `--port=2222` or `--port 2222`

### Config files

If we were to run legion with the following parameters

```bash
$ ./legion --host localhost \
           --port 2222 \
           --username root \
           --password fredfred \
           --set host pi1.local \
           --set timezone Europe/London \
           --set admin fred \
           scripts/bootstrap.legion
```

We could create a config file like this (called `server.txt`):

```
--host localhost
--port 2222
--username root
--password fredfred
--set host pi1.local
--set timezone Europe/London
--set admin fred
```

Note the lack of the line continuation markers (`\`). Now we are able to run legion like this:

```bash
$ ./legion --config server.txt scripts/bootstrap.legion
```

Which simplifies things when we have multiple scripts to run. Command line parameters are processed in order so `--config server.txt --username admin` will overwrite the `username` value from `server.txt` (`root`) with the value `admin`

## Legion commands

Legion has very few command

### `#` - A comment

Any line starting with a `#` is treated as a comment

### `ex` - Execute a command

The line `ex ls -l` will run the command `ls -l` on the target machine. The output of which will appear in the log file and on the screen

### `copy` - Copy a file from the host to the target

The command `copy fred.txt albert.txt` will copy the file `fred.txt` from the host machine to `albert.txt` on the target machine. If the target file exists and is writeable it will overwrite the file, if not it will error

### `run` - Run a script from the host on the target

The command `run create_admin_group.sh` will copy the file `create_admin_group.sh` from the host machine to the target machine, make it executable and run it. Any output from the script will appear in the log file and on the screen

Once the script has been run it will be removed from the target machine

### `call` - Run another legion scripts

Rather than have large, do it all, files you can call other legion scripts with `call configure_filewall.legion`. Once the called script terminates execution of the current script will continue

### `set` - Set a variable to be used in templates

Legion has a mini templating system and the `set` command is used to set variables in the same was at the command line arguments

## Templating

The templating is only available within legion scripts. With our previous example we set the variable `timezone` in the config file. We can use it in a script as follows:

```
...
run set_timezone.sh {timezone}
...
```

Before this line is run `{timezone}` will be replaced by `Europe/London`

## A worked example

Would be nice

## Why the long ass file extension?

All the shorter ones were already taken. Anything shorter than the full name would have been needlessly cryptic so I went the whole hog. This way Github wont think that half this project is written in Lisp
