#!/bin/bash

psql -c "CREATE USER overseer WITH PASSWORD 'dev';"
psql -c "CREATE DATABASE overseer;"
psql -c "GRANT ALL PRIVILEGES ON DATABASE overseer to overseer;"
psql -c "ALTER USER overseer WITH SUPERUSER;"

cd base_schema
PGPASSWORD=dev psql -d overseer -U overseer -f db_setup.sql
cd ..
