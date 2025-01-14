#!/bin/bash
# If errors during the execution, exit the bash script.
set -e

# Database setup
echo "Creating greenlight database as superuser..."
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" <<-EOSQL
	CREATE DATABASE greenlight;
EOSQL

echo "Opening new database as superuser"
echo "Adding extension and creating a normal operational user"
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "greenlight" <<-EOSQL
	CREATE EXTENSION IF NOT EXISTS citext;
	CREATE ROLE greenlight WITH LOGIN PASSWORD 'secret_password';
EOSQL
