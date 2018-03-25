DROP DATABASE IF EXISTS assess;

CREATE DATABASE assess;

\c assess

CREATE TABLE assessment (
  account integer PRIMARY KEY,
  address text,
  value integer NOT NULL
);

-- Using a local version of psql, connect to the database and
-- run the following command:
--
-- \copy assessment FROM './simplified.csv' WITH (FORMAT csv);
