#!/bin/sh
set -e

export DATABASE_URL="postgresql://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME"
atlas migrate apply --url $DATABASE_URL --dir file://migrations

exec ./http
