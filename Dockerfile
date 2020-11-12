FROM alpine
LABEL maintainer = "Webank CTB Team"

ENV BASE_HOME=/app/notification
ENV NOTIFICATION_PORT=9095
ENV MAIL_AUTH_USER=default_user
ENV MAIL_AUTH_PASSWORD=default_password
ENV MAIL_AUTH_SERVER=default_server

RUN mkdir -p $BASE_HOME $BASE_HOME/conf $BASE_HOME/logs

ADD build/start.sh $BASE_HOME/
ADD build/stop.sh $BASE_HOME/
ADD build/default.json $BASE_HOME/conf/
ADD wecube-plugins-notifications $BASE_HOME/

RUN chmod +x $BASE_HOME/*.sh
RUN chmod +x $BASE_HOME/wecube-plugins-notifications

WORKDIR $BASE_HOME

ENTRYPOINT ["/bin/sh", "start.sh"]