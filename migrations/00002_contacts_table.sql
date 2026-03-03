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

CREATE TRIGGER contacts_set_updated_at_trigger
    BEFORE UPDATE
    ON contacts
    FOR EACH ROW
EXECUTE PROCEDURE contacts_set_updated_at();

COMMENT ON TABLE contacts IS 'Contact objects used by EPP (registrant/admin/tech/billing roles).';
COMMENT ON COLUMN contacts.id IS 'Internal surrogate key.';
COMMENT ON COLUMN contacts.contact_id IS 'Registry contact identifier visible to EPP clients.';
COMMENT ON COLUMN contacts.roid IS 'Repository object identifier for the contact object.';
COMMENT ON COLUMN contacts.registrar_id IS 'Current sponsoring registrar.';
COMMENT ON COLUMN contacts.name IS 'Primary contact name.';
COMMENT ON COLUMN contacts.email IS 'Primary contact email address.';
COMMENT ON COLUMN contacts.organization IS 'Organization name, if provided.';
COMMENT ON COLUMN contacts.phone IS 'Voice phone number in E.164-compatible format.';
COMMENT ON COLUMN contacts.fax IS 'Fax number, if provided.';
COMMENT ON COLUMN contacts.auth_info_hash IS 'Hashed authInfo secret used for sensitive contact operations.';
COMMENT ON COLUMN contacts.disclosure_flags IS 'Disclosure preferences for contact data publication.';
COMMENT ON COLUMN contacts.created_at IS 'Row creation timestamp.';
COMMENT ON COLUMN contacts.updated_at IS 'Row last update timestamp.';
COMMENT ON COLUMN contacts.created_by_client_id IS 'EPP client identifier that created the contact (crID).';
COMMENT ON COLUMN contacts.updated_by_client_id IS 'EPP client identifier that last updated the contact (upID).';
COMMENT ON COLUMN contacts.deleted_at IS 'Soft-delete timestamp; NULL means active.';
COMMENT ON COLUMN contacts.postal_name IS 'Postal info contact name.';
COMMENT ON COLUMN contacts.postal_org IS 'Postal info organization.';
COMMENT ON COLUMN contacts.postal_code IS 'Postal/ZIP code.';
COMMENT ON COLUMN contacts.city IS 'City/locality.';
COMMENT ON COLUMN contacts.country_code IS 'ISO 3166-1 alpha-2 country code.';
COMMENT ON COLUMN contacts.street1 IS 'Postal street line 1.';
COMMENT ON COLUMN contacts.street2 IS 'Postal street line 2.';
COMMENT ON COLUMN contacts.street3 IS 'Postal street line 3.';
COMMENT ON COLUMN contacts.state_province IS 'State or province.';

-- +goose Down
DROP INDEX IF EXISTS contacts_contact_id_idx;
DROP INDEX IF EXISTS contacts_registrar_id_idx;
DROP INDEX IF EXISTS contacts_email_idx;
DROP INDEX IF EXISTS contacts_deleted_at_idx;
DROP TRIGGER IF EXISTS contacts_set_updated_at_trigger ON contacts;
DROP FUNCTION IF EXISTS contacts_set_updated_at();
DROP TABLE IF EXISTS contacts;
