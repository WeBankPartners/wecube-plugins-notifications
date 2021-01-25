#!/bin/bash

echo "start run"
sed -i "s~{{MAIL_AUTH_USER}}~$MAIL_AUTH_USER~g" conf/default.json
sed -i "s~{{MAIL_AUTH_PASSWORD}}~$MAIL_AUTH_PASSWORD~g" conf/default.json
sed -i "s~{{MAIL_AUTH_SERVER}}~$MAIL_AUTH_SERVER~g" conf/default.json


if [ -n "$NOTIFICATION_LOCAL_DNS_MAP" ]
then
  dns_map=${NOTIFICATION_LOCAL_DNS_MAP}
  dns_map_split=(${dns_map//,/ })
  for v in ${dns_map_split[@]}
  do
    echo ${v//=/ } >> /etc/hosts
  done
fi

./wecube-plugins-notifications >> logs/app.log 2>&1


