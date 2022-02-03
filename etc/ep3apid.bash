#!/bin/bash

PATH="/usr/local/bin:/usr/bin:/bin"
export PATH

#
# Setup environment
#
cd "$(dirname "$0")/.."
echo "Working directory $(pwd)"

#
# Start the webserver with an environment from /Sites/SITENAME/etc/cait.bash
#

function startService() {
        echo "Starting ep3apid"
        /usr/local/bin/ep3apid /usr/local/etc/eprinttools/settings.json
}

function stopService() {
        echo "Stopping ep3apid"
        for PID in $(pgrep ep3apid); do
                kill -s TERM "$PID"
        done
}

function statusService() {
        for PID in $(pgrep ep3apid); do
                echo "CAIT running as $PID"
        done
}

# Handle requested action
case "$1" in
        start)
                startService
                ;;
        stop)
                stopService
                ;;
        restart)
                stopService
                startService
                ;;
        status)
                statusService
                ;;
        *)
                echo 'usage: epi3apid [start|stop|restart|reindex|status]'
                ;;
esac



