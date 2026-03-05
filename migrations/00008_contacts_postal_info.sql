-- +goose Up
CREATE TYPE contacts_postal_info_type AS ENUM ('loc', 'int');

CREATE TABLE IF NOT EXISTS contacts_postal_info
(
    id             BIGSERIAL PRIMARY KEY,
    contact_id     BIGINT                                NOT NULL,
    name           TEXT        DEFAULT NULL,
    type           contacts_postal_info_type             NOT NULL,
    postal_name    TEXT        DEFAULT NULL,
    postal_org     TEXT        DEFAULT NULL,
    postal_code    TEXT        DEFAULT NULL,
    city           TEXT        DEFAULT NULL,
    country_code   TEXT        DEFAULT NULL,
    streets        JSONB       DEFAULT '[]'::JSONB,
    state_province TEXT        DEFAULT NULL,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT fk_contact_id FOREIGN KEY (contact_id) REFERENCES contacts (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS contacts_contact_id_type_uidx ON contacts_postal_info (contact_id, type);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION contacts_postal_info_set_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER contacts_postal_info_set_updated_at_trigger
    BEFORE UPDATE
    ON contacts_postal_info
    FOR EACH ROW
EXECUTE PROCEDURE contacts_postal_info_set_updated_at();

COMMENT ON TABLE contacts_postal_info IS 'Postal info for contacts.';
COMMENT ON COLUMN contacts_postal_info.id IS 'Internal surrogate key.';
COMMENT ON COLUMN contacts_postal_info.contact_id IS 'Contact ID.';
COMMENT ON COLUMN contacts_postal_info.name IS 'Contact name.';
COMMENT ON COLUMN contacts_postal_info.type IS 'Type of postal info.';
COMMENT ON COLUMN contacts_postal_info.postal_name IS 'Postal info contact name.';
COMMENT ON COLUMN contacts_postal_info.postal_org IS 'Postal info organization.';
COMMENT ON COLUMN contacts_postal_info.postal_code IS 'Postal/ZIP code.';
COMMENT ON COLUMN contacts_postal_info.city IS 'City/locality.';
COMMENT ON COLUMN contacts_postal_info.country_code IS 'ISO 3166-1 alpha-2 country code.';
COMMENT ON COLUMN contacts_postal_info.streets IS 'List of streets.';
COMMENT ON COLUMN contacts_postal_info.state_province IS 'State or province.';

-- +goose Down
DROP TRIGGER IF EXISTS contacts_postal_info_set_updated_at_trigger ON contacts_postal_info;
DROP FUNCTION IF EXISTS contacts_postal_info_set_updated_at();
DROP INDEX IF EXISTS contacts_contact_id_type_uidx;
DROP TABLE IF EXISTS contacts_postal_info;
DROP TYPE IF EXISTS contacts_postal_info_type;
