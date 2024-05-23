# Legion - scriptable access to remote servers via SSH

Legion is a tool that allows a simple scripting language (with a mini template system) to run over SSH to a target system. The scripting language can issue commands, copy files from the host to the target and upload and execute stored scripts. I use it to automate the configuration of servers. It probably has other uses and there are certainly a million and one other (better) tools that do the same thing

Why is it called "Legion"? Well I have created several such system before of such complexity that it made it harder to deploy despite their sweet architecture. I've written so many of them they are legion :)

This time I have gone for the simplest solution

## Using Legion

## Legion commands

Legion has very few command

### `#` - A comment

Any line starting with a `#` is treated as a comment. Blank lines are skipped

### `CMD` - Execute a command

The line `CMD ls -l` will run the command `ls -l` on the target machine. The output of which will appear in the log file and on the screen

### `COPY` - Copy a file from the host to the target

The command `COPY fred.txt albert.txt` will copy the file `fred.txt` from the host machine to `albert.txt` on the target machine. If the target file exists and is writeable it will overwrite the file, if not it will error

### `CONFIG` - Set an internal variable

To set a variable `CONFIG timezone Europe/London` which can be used later in templates

### `ECHO` - Write a message out to the log

The line `ECHO Install Nginx` will display "Install Nginx" to the console and log

### `DEBUG` - Display all the variables set with `CONFIG`

Write out all the variables that were set by `CONFIG`

### `HALT` - Halt the process

Once the `HALT` line is execute Legion will stop

### `INCLUDE` - Start processing a new legion file

Once `INCLUDE install_postgres.legion` is encountered the current script is suspended and `install_postgres.legion` will be run. Once `install_postgres.legion` completed the original script will continue to run

## Templating

The templating is only available within legion scripts. With our previous example we set the variable `timezone` in the config file. We can use it in a script as follows:

```
CMD set_timezone.sh {{timezone}}
```

Before this line is run `{{timezone}}` will be replaced by `Europe/London`

## A worked example

Would be nice

## Why the long ass file extension?

All the shorter ones were already taken. Anything shorter than the full name would have been needlessly cryptic so I went the whole hog. This way Github wont think that half this project is written in Lisp
