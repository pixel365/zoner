-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION zones_available_periods_is_valid(periods INT2[]) RETURNS BOOLEAN AS
$$
SELECT
    cardinality(periods) > 0
    AND NOT EXISTS (
        SELECT 1
        FROM unnest(periods) AS p
        WHERE p IS NULL OR p <= 0
    )
    AND (
        SELECT COUNT(*) = COUNT(DISTINCT p)
        FROM unnest(periods) AS p
    );
$$ LANGUAGE SQL IMMUTABLE;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS zones
(
    id                             BIGSERIAL PRIMARY KEY,
    name                           TEXT                                  NOT NULL,
    is_active                      BOOLEAN                               NOT NULL DEFAULT TRUE,
    available_periods              INT2[]                               NOT NULL DEFAULT ARRAY [1]::INT2[],
    add_grace_period_days          INT2                                 NOT NULL DEFAULT 0,
    renew_grace_period_days        INT2                                 NOT NULL DEFAULT 0,
    transfer_grace_period_days     INT2                                 NOT NULL DEFAULT 0,
    redemption_grace_period_days   INT2                                 NOT NULL DEFAULT 0,
    pending_delete_period_days     INT2                                 NOT NULL DEFAULT 0,
    restore_grace_period_days      INT2                                 NOT NULL DEFAULT 0,
    created_at                     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at                     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT zones_name_check CHECK (
        starts_with(lower(name), '.') AND char_length(name) > 1 AND char_length(name) < 63
        ),
    CONSTRAINT zones_available_periods_check CHECK (zones_available_periods_is_valid(available_periods)),
    CONSTRAINT zones_add_grace_period_days_check CHECK (add_grace_period_days >= 0),
    CONSTRAINT zones_renew_grace_period_days_check CHECK (renew_grace_period_days >= 0),
    CONSTRAINT zones_transfer_grace_period_days_check CHECK (transfer_grace_period_days >= 0),
    CONSTRAINT zones_redemption_grace_period_days_check CHECK (redemption_grace_period_days >= 0),
    CONSTRAINT zones_pending_delete_period_days_check CHECK (pending_delete_period_days >= 0),
    CONSTRAINT zones_restore_grace_period_days_check CHECK (restore_grace_period_days >= 0)
);

CREATE UNIQUE INDEX zones_name_uidx ON zones (lower(name));
CREATE INDEX zones_is_active_idx ON zones (is_active);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION zones_set_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION zones_sort_available_periods() RETURNS TRIGGER AS
$$
BEGIN
    NEW.available_periods = ARRAY(
        SELECT p
        FROM unnest(NEW.available_periods) AS p
        ORDER BY p
    );
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

CREATE TRIGGER zones_sort_available_periods_trigger
    BEFORE INSERT OR UPDATE
    ON zones
    FOR EACH ROW
EXECUTE PROCEDURE zones_sort_available_periods();

CREATE TRIGGER zones_set_updated_at_trigger
    BEFORE UPDATE
    ON zones
    FOR EACH ROW
EXECUTE PROCEDURE zones_set_updated_at();

COMMENT ON TABLE zones IS 'Supported top-level zones and their registration lifecycle policy.';
COMMENT ON COLUMN zones.id IS 'Internal surrogate key.';
COMMENT ON COLUMN zones.name IS 'Zone name served by registry, usually with leading dot (for example, .com).';
COMMENT ON COLUMN zones.is_active IS 'Whether zone accepts new EPP operations.';
COMMENT ON COLUMN zones.available_periods IS 'Allowed registration periods for this zone (for example, {1,3,5}).';
COMMENT ON COLUMN zones.add_grace_period_days IS 'Add Grace Period length in days.';
COMMENT ON COLUMN zones.renew_grace_period_days IS 'Renew Grace Period length in days.';
COMMENT ON COLUMN zones.transfer_grace_period_days IS 'Transfer Grace Period length in days.';
COMMENT ON COLUMN zones.redemption_grace_period_days IS 'Redemption Grace Period length in days.';
COMMENT ON COLUMN zones.pending_delete_period_days IS 'Pending Delete period length in days.';
COMMENT ON COLUMN zones.restore_grace_period_days IS 'Restore Grace Period length in days.';
COMMENT ON COLUMN zones.created_at IS 'Row creation timestamp.';
COMMENT ON COLUMN zones.updated_at IS 'Row last update timestamp.';

-- +goose Down
DROP TRIGGER IF EXISTS zones_sort_available_periods_trigger ON zones;
DROP TRIGGER IF EXISTS zones_set_updated_at_trigger ON zones;
DROP FUNCTION IF EXISTS zones_sort_available_periods();
DROP FUNCTION IF EXISTS zones_set_updated_at();
DROP INDEX IF EXISTS zones_is_active_idx;
DROP INDEX IF EXISTS zones_name_uidx;
DROP TABLE IF EXISTS zones;
DROP FUNCTION IF EXISTS zones_available_periods_is_valid(INT2[]);
