#!/bin/sh

# Get a list of all the postgres databases and back them up to BACKUP_ROOT
# Also keep LIMIT backups, purging the oldest backups

BACKUP_ROOT=/var/local/backup

DATE=`date "+%Y-%m-%dT%H-%M-%S"`

LIMIT=15

message()
{
  printf "%s %s [%5s] %s\n" `date "+%Y-%m-%d %H:%M:%S"` $$ "$@" >> $BACKUP_ROOT/db_backup.log 2>&1
}

purge()
{
  local DB=$1

  message "Looking to purge old files for $DB";
  PATTERN="$BACKUP_ROOT/$DB/$DB.????-??-??T??-??-??.backup";

  COUNTER=0
  for X in `ls -1 $PATTERN | sort -r`
  do
    COUNTER=$((COUNTER + 1))
    if [ $COUNTER -gt $LIMIT ]
    then
      message "Deleting $X"
      rm $X
    else
      message "Keeping $X"
    fi
  done

  message "Purge of $DB complete"
}

message "Starting to process the backups"

for DB in `psql -l | grep '|' | cut -d'|' -f 1`
do
  case $DB in
    "Name")
      # This is the column header
      ;;
    "template0"|"template1"|"postgres")
      # We ignore these internal files
      ;;
    *test)
      # Ignore the test files
      ;;
    *)
      if [ ! -d "$BACKUP_ROOT/$DB" ]
      then
        mkdir -p $BACKUP_ROOT/$DB
      fi
      message "Processing $DB";
      message "Backup to $BACKUP_ROOT/$DB/$DB.$DATE.backup";
      pg_dump -F c -f $BACKUP_ROOT/$DB/$DB.$DATE.backup $DB;
      purge $DB
      ;;
  esac
done

message "Finished processing the backups"
