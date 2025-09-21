#!/usr/bin/sh

cp -v ./deploy/docker/develop/consumer/supervisor.conf /etc/supervisor/conf.d/consumer.conf

/usr/bin/supervisord -c /etc/supervisor/supervisord.conf