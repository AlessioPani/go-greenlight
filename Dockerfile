# Init script for Postgres and update permissions.
FROM postgres

COPY ./db-data/init.sh /docker-entrypoint-initdb.d/

RUN chmod 755 /docker-entrypoint-initdb.d/init.sh
