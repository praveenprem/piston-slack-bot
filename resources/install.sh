#!/usr/bin/env bash

SERVICE_DIR=/etc/systemd/system/
APP_HOME=/var/testbed-bot/
NGINX_DIR=/etc/nginx/sites-enabled/

mkdir -p $APP_HOME

if [ -f config.json ]; then
  mv config.json $APP_HOME
fi

if [ -f testbed-bot ]; then
  mv testbed-bot $APP_HOME
fi

if [[ $(systemctl --no-pager list-unit-files | grep testbed-bot) == 0 ]]; then
  systemctl restart testbed-bot.service
else
  mv testbed-bot.service $SERVICE_DIR
  systemctl daemon-reload
  systemctl enable testbed-bot.service
  systemctl start testbed-bot.service
fi

if [ ! -f $NGINX_DIR/piston.razorcorp.dev.nginx ]; then
  mv piston.razorcorp.dev.nginx $NGINX_DIR
  systemctl restart nginx.service
fi
