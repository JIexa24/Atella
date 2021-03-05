#! /usr/bin/env bash

# chkconfig: 2345 99 01
# description: Atella daemon

### BEGIN INIT INFO
# Provides:          atella
# Required-Start:    $all
# Required-Stop:     $remote_fs $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Start atella at boot time
### END INIT INFO

# this init script supports three different variations:
#  1. New lsb that define start-stop-daemon
#  2. Old lsb that don't have start-stop-daemon but define, log, pidofproc and killproc
#  3. Centos installations without lsb-core installed
#
# In the third case we have to define our own functions which are very dumb
# and expect the args to be positioned correctly.

# Command-line options that can be set in /etc/default/atella.  These will override
# any config file values.
ATELLA_OPTS=

USER=atella
GROUP=atella

DEFAULT=/etc/default/atella

# Process name ( For display )
name=atella

# Daemon name, where is the actual executable
daemon=/usr/bin/atella
daemon_cli=/usr/bin/atella-cli

# If the daemon is not there, then exit.
[ -x $daemon ] || exit 5
# If the daemon-cli is not there, then exit.
[ -x $daemon ] || exit 6

# pid file for the daemon
pidfile=`/usr/bin/atella-cli -print-pidfile 2>/dev/null`
piddir=`dirname $pidfile`

if [ ! -d "$piddir" ]; then
    mkdir -p $piddir
    chown $USER:$GROUP $piddir
fi

# Configuration file
config=/etc/atella/atella.conf
confdir=/etc/atella/conf.d

case $1 in
    start)
        ;;

    stop)
        ;;

    reload)
        ;;

    restart)
        $0 stop && sleep 2 && $0 start
        ;;

    status)
        ;;

    version)
        $daemon -version
        $daemon_cli -version
        ;;

    *)
        # For invalid arguments, print the usage message.
        echo "Usage: $0 {start|stop|restart|status|version}"
        exit 2
        ;;
esac
