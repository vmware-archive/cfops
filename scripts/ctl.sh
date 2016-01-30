#!/bin/bash

RUN_DIR='/root/pcfbackup'
LOG_DIR='/root/pcfbackup/log'
PIDFILE="$RUN_DIR/pid"

CRON_STATUS=`pgrep cron | wc -l`
CRON_ENTRY=$(cat /root/scripts/backup.cron)
CRONTAB_TMP="$RUN_DIR/crontab.$$.tmp"

case $1 in

  start)

    mkdir -p "$RUN_DIR" "$LOG_DIR"

    echo $$ > $PIDFILE

    if [ $CRON_STATUS -eq 0 ]; then
      /usr/sbin/cron start
    fi

    if crontab -l | sed -e 's/^#.*//' | grep "$CRON_ENTRY" 2>&1>/dev/null; then
      echo 'Already running'
    else
      crontab -l | grep -v "$CRON_ENTRY" > "$CRONTAB_TMP"
      echo "$CRON_ENTRY" >> "$CRONTAB_TMP"
      crontab < "$CRONTAB_TMP"
      rm "$CRONTAB_TMP"
    fi

    while true; do
      sleep 60
    done

    ;;

  stop)

    if crontab -l | sed -e 's/^#.*//' | grep "$CRON_ENTRY" 2>&1>/dev/null; then
      crontab -l | grep -v "$CRON_ENTRY" > "$CRONTAB_TMP"
      crontab < "$CRONTAB_TMP"
      rm "$CRONTAB_TMP"
    else
      echo 'Not running'
    fi

    kill -TERM $(cat "$PIDFILE")

    rm -f $PIDFILE

    ;;

  *)
    echo "Usage: ctl {start|stop}" ;;

esac
