#!/bin/bash
# If errors during the execution, exit the bash script.
set -e

# Database setup
echo "Adding extension and creating a normal operational user"
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "greenlight" <<-EOSQL
	CREATE EXTENSION IF NOT EXISTS citext;
	CREATE ROLE greenlight WITH LOGIN PASSWORD 'secret_password';
	ALTER DATABASE greenlight OWNER TO greenlight;
EOSQL
