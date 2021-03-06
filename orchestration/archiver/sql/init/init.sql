-- kubectl port-forward svc/mypostgres 5436:5432 -n postgres
-- make init
-- user creation
CREATE USER IF NOT EXISTS archiver_app WITH PASSWORD 'DB_PASSWORD';

CREATE DATABASE IF NOT EXISTS  archiverDB_SUFFIX;

GRANT ALL PRIVILEGES ON DATABASE archiverDB_SUFFIX TO archiver_app;