BEGIN;

DROP TABLE IF EXISTS outbox;

DROP TYPE outbox_event;

DROP TYPE outbox_status;

COMMIT;
