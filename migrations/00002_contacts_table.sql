-- +goose Up
CREATE TABLE IF NOT EXISTS contacts
(
    id                   BIGSERIAL PRIMARY KEY,
    contact_id           TEXT                                  NOT NULL,
    roid                 TEXT                                  NOT NULL UNIQUE,
    registrar_id         BIGINT                                NOT NULL,
    name                 TEXT                                  NOT NULL,
    email                TEXT                                  NOT NULL,
    organization         TEXT,
    phone                TEXT,
    fax                  TEXT,
    auth_info_hash       TEXT                                  NOT NULL,
    disclosure_flags     JSONB                                 NOT NULL DEFAULT '{}'::JSONB,
    created_at           TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at           TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by_client_id TEXT                                  NOT NULL,
    updated_by_client_id TEXT,
    deleted_at           TIMESTAMPTZ,

    postal_name          TEXT                                  NOT NULL,
    postal_org           TEXT,
    postal_code          TEXT,
    city                 TEXT                                  NOT NULL,
    country_code         TEXT                                  NOT NULL,
    street1              TEXT,
    street2              TEXT,
    street3              TEXT,
    state_province       TEXT,

    CONSTRAINT fk_registrar_id FOREIGN KEY (registrar_id) REFERENCES registrars (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS contacts_contact_id_idx ON contacts (lower(contact_id));
CREATE INDEX IF NOT EXISTS contacts_registrar_id_idx ON contacts (registrar_id);
CREATE INDEX IF NOT EXISTS contacts_email_idx ON contacts (lower(email));
CREATE INDEX IF NOT EXISTS contacts_deleted_at_idx ON contacts (deleted_at);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION contacts_set_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER contacts_set_updated_at
    BEFORE UPDATE
    ON contacts
    FOR EACH ROW
EXECUTE PROCEDURE contacts_set_updated_at();

-- +goose Down
DROP INDEX IF EXISTS contacts_contact_id_idx;
DROP INDEX IF EXISTS contacts_registrar_id_idx;
DROP INDEX IF EXISTS contacts_email_idx;
DROP INDEX IF EXISTS contacts_deleted_at_idx;
DROP TRIGGER IF EXISTS contacts_set_updated_at ON contacts;
DROP FUNCTION IF EXISTS contacts_set_updated_at();
DROP TABLE IF EXISTS contacts;
