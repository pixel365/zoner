-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION contacts_disclose_is_valid(d jsonb) RETURNS boolean AS
$$
SELECT d = '{}'::jsonb
           OR (
           jsonb_typeof(d) = 'object'
               AND d ?& array ['flag','fields']
               AND (d - 'flag' - 'fields') = '{}'::jsonb
               AND jsonb_typeof(d -> 'flag') = 'number'
               AND (d ->> 'flag') IN ('0', '1')
               AND jsonb_typeof(d -> 'fields') = 'array'
               AND jsonb_array_length(d -> 'fields') > 0
               AND NOT EXISTS (SELECT 1
                               FROM jsonb_array_elements(d -> 'fields') AS e(v)
                               WHERE jsonb_typeof(v) <> 'string'
                                  OR (v #>> '{}') NOT IN (
                                                          'name', 'org', 'addr', 'voice', 'fax', 'email',
                                                          'addr:int', 'addr:loc'
                                   ))
               AND (SELECT count(*) = count(DISTINCT (v #>> '{}'))
                    FROM jsonb_array_elements(d -> 'fields') AS e(v))
           );
$$
    LANGUAGE SQL IMMUTABLE;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS contacts
(
    id                   BIGSERIAL PRIMARY KEY,
    contact_id           TEXT                                  NOT NULL,
    roid                 TEXT                                  NOT NULL,
    registrar_id         BIGINT                                NOT NULL,
    name                 TEXT                                  NOT NULL,
    email                TEXT                                  NOT NULL,
    organization         TEXT        DEFAULT NULL,
    voice                TEXT        DEFAULT NULL,
    fax                  TEXT        DEFAULT NULL,
    auth_info_hash       TEXT                                  NOT NULL,
    disclose             JSONB                                 NOT NULL DEFAULT '{}'::JSONB,
    created_at           TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at           TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by_client_id TEXT        DEFAULT NULL,
    updated_by_client_id TEXT        DEFAULT NULL,
    deleted_at           TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT fk_registrar_id FOREIGN KEY (registrar_id) REFERENCES registrars (id),
    CONSTRAINT contacts_disclose_check CHECK (contacts_disclose_is_valid(disclose))
);

CREATE UNIQUE INDEX IF NOT EXISTS contacts_contact_id_idx ON contacts (lower(contact_id));
CREATE INDEX IF NOT EXISTS contacts_registrar_id_idx ON contacts (registrar_id);
CREATE INDEX IF NOT EXISTS contacts_email_idx ON contacts (lower(email));
CREATE INDEX IF NOT EXISTS contacts_deleted_at_idx ON contacts (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS contacts_roid_idx ON contacts (upper(roid));

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION contacts_set_roid() RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.id IS NULL THEN
        NEW.id := nextval(pg_get_serial_sequence('contacts', 'id'));
    END IF;

    NEW.roid = 'C' || NEW.id::text || '-REG';
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION contacts_set_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.roid <> OLD.roid THEN
        RAISE EXCEPTION 'Cannot change roid of contact';
    END IF;

    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER contacts_set_roid_trigger
    BEFORE INSERT
    ON contacts
    FOR EACH ROW
EXECUTE PROCEDURE contacts_set_roid();

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
COMMENT ON COLUMN contacts.voice IS 'Voice phone number in E.164-compatible format.';
COMMENT ON COLUMN contacts.fax IS 'Fax number, if provided.';
COMMENT ON COLUMN contacts.auth_info_hash IS 'Hashed authInfo secret used for sensitive contact operations.';
COMMENT ON COLUMN contacts.disclose IS 'Disclosure preferences for contact data publication.';
COMMENT ON COLUMN contacts.created_at IS 'Row creation timestamp.';
COMMENT ON COLUMN contacts.updated_at IS 'Row last update timestamp.';
COMMENT ON COLUMN contacts.created_by_client_id IS 'EPP client identifier that created the contact (crID).';
COMMENT ON COLUMN contacts.updated_by_client_id IS 'EPP client identifier that last updated the contact (upID).';
COMMENT ON COLUMN contacts.deleted_at IS 'Soft-delete timestamp; NULL means active.';

-- +goose Down
DROP INDEX IF EXISTS contacts_contact_id_idx;
DROP INDEX IF EXISTS contacts_registrar_id_idx;
DROP INDEX IF EXISTS contacts_email_idx;
DROP INDEX IF EXISTS contacts_deleted_at_idx;
DROP INDEX IF EXISTS contacts_roid_idx;
DROP TRIGGER IF EXISTS contacts_set_roid_trigger ON contacts;
DROP TRIGGER IF EXISTS contacts_set_updated_at_trigger ON contacts;
DROP FUNCTION IF EXISTS contacts_set_roid();
DROP FUNCTION IF EXISTS contacts_set_updated_at();
DROP TABLE IF EXISTS contacts;
DROP FUNCTION IF EXISTS contacts_disclose_is_valid(jsonb);
