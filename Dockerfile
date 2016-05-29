FROM busybox

COPY bin/mysql-statsd /
COPY script/run.sh /

ENTRYPOINT ["/run.sh"]
