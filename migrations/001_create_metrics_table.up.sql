-- Generated with `migrate create -ext sql -dir migrations -seq -digits 3 create_metrics_table`

BEGIN;

CREATE TYPE metric_type AS ENUM ('counter', 'gauge');

CREATE TABLE IF NOT EXISTS metrics (
    id      serial PRIMARY KEY,
    name    varchar(64) UNIQUE NOT NULL,
    type    metric_type        NOT NULL,
    counter bigint             NULL,
    gauge   double precision   NULL
);

COMMIT;
