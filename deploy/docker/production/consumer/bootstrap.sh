#!/usr/bin/sh

cp -v ./deploy/docker/production/consumer/supervisor.conf /etc/supervisor/conf.d/consumer.conf

/usr/bin/supervisord -c /etc/supervisor/supervisord.conf