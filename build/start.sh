#!/bin/bash

echo "start run"
sed -i "s~{{MAIL_AUTH_USER}}~$MAIL_AUTH_USER~g" conf/default.json
sed -i "s~{{MAIL_AUTH_PASSWORD}}~$MAIL_AUTH_PASSWORD~g" conf/default.json
sed -i "s~{{MAIL_AUTH_SERVER}}~$MAIL_AUTH_SERVER~g" conf/default.json

./wecube-plugins-notifications >> logs/app.log 2>&1


