FROM debian:jessie

COPY bin/mysql-statsd /
COPY script/run.sh /

ENTRYPOINT ["/run.sh"]
