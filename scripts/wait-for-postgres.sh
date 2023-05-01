#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
    >&2 echo "Waiting for PostgreSQL"
    sleep 1
done

>&2 echo "PostgreSQL is up"
exec $cmd