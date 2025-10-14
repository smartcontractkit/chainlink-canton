#!/usr/bin/env bash

set -Eeo pipefail

function create_database() {
    local database=$1
    echo "    Creating database: '$database'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        CREATE DATABASE "$database";
        GRANT ALL PRIVILEGES ON DATABASE "$database" TO $POSTGRES_USER;
EOSQL
}

if [ -n "$POSTGRES_INIT_DATABASES" ]; then
    echo "Creating multiple databases: $POSTGRES_INIT_DATABASES"
    for database in $(echo $POSTGRES_INIT_DATABASES | tr ',' ' '); do
        create_database $database
    done
    echo "All databases created"
fi