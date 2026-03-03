-- +goose Up
CREATE TABLE IF NOT EXISTS domains
(
    id                    BIGSERIAL PRIMARY KEY,
    name                  TEXT        NOT NULL,
    punycode              TEXT        NOT NULL,
    roid                  TEXT        NOT NULL UNIQUE,
    registrar_id          BIGINT      NOT NULL,
    registrant_contact_id BIGINT      NOT NULL,
    auth_info_hash        TEXT        NOT NULL,
    period_value          INT8        NOT NULL,
    period_unit           TEXT        NOT NULL DEFAULT 'y',
    created_at            TIMESTAMPTZ          DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at            TIMESTAMPTZ          DEFAULT CURRENT_TIMESTAMP NOT NULL,
    registered_at         TIMESTAMPTZ NOT NULL,
    expires_at            TIMESTAMPTZ NOT NULL,
    transferred_at        TIMESTAMPTZ,
    deleted_at            TIMESTAMPTZ,
    is_locked             BOOLEAN     NOT NULL DEFAULT FALSE,
    is_auto_renew         BOOLEAN     NOT NULL DEFAULT TRUE,

    CONSTRAINT domains_registered_at_check CHECK (registered_at <= expires_at),
    CONSTRAINT domains_period_value_check CHECK (period_value > 0),
    CONSTRAINT domains_period_unit_check CHECK (period_unit = 'y'),
    CONSTRAINT fk_registrar_id FOREIGN KEY (registrar_id) REFERENCES registrars (id),
    CONSTRAINT fk_registrant_contact_id FOREIGN KEY (registrant_contact_id) REFERENCES contacts (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS domains_punycode_uidx ON domains (lower(punycode));
CREATE INDEX IF NOT EXISTS domains_registrar_id_idx ON domains (registrar_id);
CREATE INDEX IF NOT EXISTS domains_registrant_contact_id_idx ON domains (registrant_contact_id);
CREATE INDEX IF NOT EXISTS domains_expires_at_idx ON domains (expires_at);
CREATE INDEX IF NOT EXISTS domains_deleted_at_idx ON domains (deleted_at);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION domains_set_registered_at() RETURNS TRIGGER AS
$$
BEGIN
    NEW.registered_at = CURRENT_TIMESTAMP;
    NEW.expires_at = NEW.registered_at + NEW.period_value * interval '1 ' || NEW.period_unit;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION domain_set_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

CREATE TRIGGER domains_set_registered_at
    BEFORE INSERT
    ON domains
    FOR EACH ROW
EXECUTE PROCEDURE domains_set_registered_at();

CREATE TRIGGER domain_set_updated_at
    BEFORE UPDATE
    ON domains
    FOR EACH ROW
EXECUTE PROCEDURE domain_set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS domain_set_updated_at ON domains;
DROP TRIGGER IF EXISTS domains_set_registered_at ON domains;
DROP FUNCTION IF EXISTS domain_set_updated_at();
DROP FUNCTION IF EXISTS domains_set_registered_at();
DROP INDEX IF EXISTS domains_expires_at_idx;
DROP INDEX IF EXISTS domains_deleted_at_idx;
DROP INDEX IF EXISTS domains_registrant_contact_id_idx;
DROP INDEX IF EXISTS domains_registrar_id_idx;
DROP INDEX IF EXISTS domains_punycode_uidx;
DROP TABLE IF EXISTS domains;
