-- +goose Up
CREATE TYPE contact_status AS ENUM (
    'ok',
    'linked',
    'clientDeleteProhibited',
    'serverDeleteProhibited',
    'clientUpdateProhibited',
    'serverUpdateProhibited',
    'pendingCreate',
    'pendingUpdate',
    'pendingDelete',
    'pendingTransfer'
    );

CREATE TYPE contact_status_source AS ENUM ('server', 'client');

CREATE TABLE contact_statuses
(
    id                BIGSERIAL PRIMARY KEY,
    contact_id        BIGINT                NOT NULL,
    status            contact_status        NOT NULL,
    source            contact_status_source NOT NULL DEFAULT contact_status_source 'server',
    created_at        TIMESTAMPTZ                    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by_client TEXT,
    reason            TEXT,

    CONSTRAINT fk_contact_id FOREIGN KEY (contact_id) REFERENCES contacts (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS contact_status_uidx ON contact_statuses (contact_id, status);
CREATE INDEX IF NOT EXISTS contact_status_idx ON contact_statuses (status);
CREATE INDEX IF NOT EXISTS contact_contact_id_idx ON contact_statuses (contact_id);

COMMENT ON TYPE contact_status IS 'Allowed EPP contact status values.';

COMMENT ON TYPE contact_status_source IS 'Origin of contact status assignment (client or server).';

COMMENT ON TABLE contact_statuses IS 'Status assignments for EPP contact objects.';
COMMENT ON COLUMN contact_statuses.id IS 'Internal surrogate key.';
COMMENT ON COLUMN contact_statuses.contact_id IS 'Reference to contacts.id.';
COMMENT ON COLUMN contact_statuses.status IS 'Assigned contact status value.';
COMMENT ON COLUMN contact_statuses.source IS 'Origin of status assignment.';
COMMENT ON COLUMN contact_statuses.created_at IS 'Status assignment timestamp.';
COMMENT ON COLUMN contact_statuses.created_by_client IS 'Client identifier that set the status, if any.';
COMMENT ON COLUMN contact_statuses.reason IS 'Optional reason for status assignment.';

-- +goose Down
DROP INDEX IF EXISTS contact_status_uidx;
DROP INDEX IF EXISTS contact_status_idx;
DROP INDEX IF EXISTS contact_contact_id_idx;
DROP TABLE IF EXISTS contact_statuses;
DROP TYPE IF EXISTS contact_status_source;
DROP TYPE IF EXISTS contact_status;
