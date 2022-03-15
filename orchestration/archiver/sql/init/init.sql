-- kubectl port-forward svc/mypostgres 5436:5432 -n postgres
-- make init
-- user creation
CREATE USER archiver_write WITH PASSWORD 'DB_PASSWORD';

CREATE USER archiver_read WITH PASSWORD 'DB_PASSWORD';

CREATE DATABASE archiverDB_SUFFIX;

GRANT ALL PRIVILEGES ON DATABASE archiverDB_SUFFIX TO archiver_write;

GRANT CONNECT ON DATABASE archiverDB_SUFFIX TO archiver_read;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO archiver_read;