#!/bin/sh
# wait-for-postgres.sh
set -e

host='$1'
shift
cmd="$@"

until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -c '\q'; do 
    >&2 echo "Postgres is inavailable - sleeping"
    sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd