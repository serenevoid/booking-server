FROM postgres:alpine

ADD scripts/init.sql /docker-entrypoint-initdb.d/
ADD scripts/load.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
