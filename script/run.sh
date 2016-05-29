#!/bin/sh

set -x

exec ./mysql-statsd -username $USERNAME -password $PASSWORD -host $MYSQL_HOST \
       -port $MYSQL_PORT -database $DATABSE -statsd_host $STATSD_SERVICE_HOST:$STATDS_SERVICE_PORT
