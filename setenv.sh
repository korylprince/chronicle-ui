#!/bin/sh

export SQLDSN="$MYSQL_USER:$MYSQL_PASS@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE?parseTime=true"

exec $1
