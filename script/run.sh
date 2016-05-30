#!/bin/sh

exec /mysql-statsd -username $USERNAME -password $PASSWORD \
       -port $MYSQL_PORT -database $DATABASE -statsd_host "statsd:8125"
