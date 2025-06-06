BEGIN;

CREATE TYPE outbox_event AS ENUM ('UserRegistered');

CREATE TYPE outbox_status AS ENUM ('pending', 'processing', 'published', 'failed');

CREATE TABLE IF NOT EXISTS outbox (
	id BIGSERIAL PRIMARY KEY,
	event outbox_event NOT NULL,
	status outbox_status DEFAULT 'pending'::outbox_status NOT NULL,
	payload JSONB NOT NULL,
	published_at TIMESTAMP,
	error TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
