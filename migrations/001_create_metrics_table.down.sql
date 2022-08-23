-- Generated with `migrate create -ext sql -dir migrations -seq -digits 3 create_metrics_table`

BEGIN;
DROP TABLE IF EXISTS metrics;
DROP TYPE IF EXISTS metric_type;
COMMIT;
