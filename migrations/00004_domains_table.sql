-- +goose Up
CREATE TABLE IF NOT EXISTS domains
(
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT        NOT NULL,
    punycode       TEXT        NOT NULL,
    roid           TEXT        NOT NULL UNIQUE,
    registrar_id   BIGINT      NOT NULL,
    auth_info_hash TEXT        NOT NULL,
    period_value   INT8        NOT NULL,
    period_unit    TEXT        NOT NULL DEFAULT 'y',
    created_at     TIMESTAMPTZ          DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at     TIMESTAMPTZ          DEFAULT CURRENT_TIMESTAMP NOT NULL,
    registered_at  TIMESTAMPTZ NOT NULL,
    expires_at     TIMESTAMPTZ NOT NULL,
    transferred_at TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,
    is_locked      BOOLEAN     NOT NULL DEFAULT FALSE,
    is_auto_renew  BOOLEAN     NOT NULL DEFAULT TRUE,

    CONSTRAINT domains_registered_at_check CHECK (registered_at <= expires_at),
    CONSTRAINT domains_period_value_check CHECK (period_value > 0),
    CONSTRAINT domains_period_unit_check CHECK (period_unit = 'y'),
    CONSTRAINT fk_registrar_id FOREIGN KEY (registrar_id) REFERENCES registrars (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS domains_punycode_uidx ON domains (lower(punycode));
CREATE INDEX IF NOT EXISTS domains_registrar_id_idx ON domains (registrar_id);
CREATE INDEX IF NOT EXISTS domains_expires_at_idx ON domains (expires_at);
CREATE INDEX IF NOT EXISTS domains_deleted_at_idx ON domains (deleted_at);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION domains_set_required_fields() RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.id IS NULL THEN
        NEW.id := nextval(pg_get_serial_sequence('domains', 'id'));
    END IF;

    NEW.roid = 'D' || NEW.id::text || '-REG';
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
    IF NEW.roid <> OLD.roid THEN
        RAISE EXCEPTION 'Cannot change roid of domain';
    END IF;

    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

CREATE TRIGGER domains_set_registered_at_trigger
    BEFORE INSERT
    ON domains
    FOR EACH ROW
EXECUTE PROCEDURE domains_set_required_fields();

CREATE TRIGGER domain_set_updated_at_trigger
    BEFORE UPDATE
    ON domains
    FOR EACH ROW
EXECUTE PROCEDURE domain_set_updated_at();

COMMENT ON TABLE domains IS 'Domain objects managed by registrars in the registry.';
COMMENT ON COLUMN domains.id IS 'Internal surrogate key.';
COMMENT ON COLUMN domains.name IS 'Human-readable domain name as received from client.';
COMMENT ON COLUMN domains.punycode IS 'Canonical ASCII domain form used for uniqueness checks.';
COMMENT ON COLUMN domains.roid IS 'Repository object identifier for the domain object.';
COMMENT ON COLUMN domains.registrar_id IS 'Current sponsoring registrar.';
COMMENT ON COLUMN domains.auth_info_hash IS 'Hashed authInfo secret used for transfer authorization.';
COMMENT ON COLUMN domains.period_value IS 'Registration period value.';
COMMENT ON COLUMN domains.period_unit IS 'Registration period unit; currently restricted to years.';
COMMENT ON COLUMN domains.created_at IS 'Row creation timestamp.';
COMMENT ON COLUMN domains.updated_at IS 'Row last update timestamp.';
COMMENT ON COLUMN domains.registered_at IS 'Initial registration timestamp.';
COMMENT ON COLUMN domains.expires_at IS 'Domain expiration timestamp.';
COMMENT ON COLUMN domains.transferred_at IS 'Timestamp of the last successful transfer.';
COMMENT ON COLUMN domains.deleted_at IS 'Soft-delete timestamp; NULL means active.';
COMMENT ON COLUMN domains.is_locked IS 'Quick lock flag to block selected mutating operations.';
COMMENT ON COLUMN domains.is_auto_renew IS 'Whether domain should be auto-renewed at expiry.';

-- +goose Down
DROP TRIGGER IF EXISTS domain_set_updated_at_trigger ON domains;
DROP TRIGGER IF EXISTS domains_set_registered_at_trigger ON domains;
DROP FUNCTION IF EXISTS domain_set_updated_at();
DROP FUNCTION IF EXISTS domains_set_required_fields();
DROP INDEX IF EXISTS domains_expires_at_idx;
DROP INDEX IF EXISTS domains_deleted_at_idx;
DROP INDEX IF EXISTS domains_registrar_id_idx;
DROP INDEX IF EXISTS domains_punycode_uidx;
DROP TABLE IF EXISTS domains;
