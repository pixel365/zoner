-- +goose Up
CREATE TYPE staff_role AS ENUM ('root', 'admin', 'manager', 'support', 'employee');

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION staff_roles_is_valid(roles staff_role[]) RETURNS BOOLEAN AS
$$
SELECT cardinality(roles) > 0
           AND NOT EXISTS (SELECT 1
                           FROM unnest(roles) AS p
                           WHERE p IS NULL)
           AND (SELECT COUNT(*) = COUNT(DISTINCT p)
                FROM unnest(roles) AS p);
$$ LANGUAGE SQL IMMUTABLE;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS staffs
(
    id            BIGSERIAL PRIMARY KEY,
    name          TEXT         NOT NULL,
    last_name     TEXT,
    email         TEXT         NOT NULL,
    password_hash TEXT         NOT NULL,
    roles         staff_role[] NOT NULL DEFAULT ARRAY ['employee']::staff_role[],
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at    TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT staffs_roles_check CHECK (staff_roles_is_valid(roles)),
    CONSTRAINT email_check CHECK (
        char_length(email) >= 5 AND lower(email) ~* '^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$'
        )
);

CREATE UNIQUE INDEX IF NOT EXISTS staffs_email_uidx ON staffs (lower(email));
CREATE INDEX IF NOT EXISTS staffs_roles_idx ON staffs USING GIN (roles);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION staffs_set_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

CREATE TRIGGER staffs_set_updated_at
    BEFORE UPDATE
    ON staffs
    FOR EACH ROW
EXECUTE PROCEDURE staffs_set_updated_at();

COMMENT ON TYPE staff_role IS 'Role values for internal staff access model.';

COMMENT ON TABLE staffs IS 'Internal staff accounts used for operational access.';
COMMENT ON COLUMN staffs.id IS 'Internal surrogate key.';
COMMENT ON COLUMN staffs.name IS 'Given name.';
COMMENT ON COLUMN staffs.last_name IS 'Family name.';
COMMENT ON COLUMN staffs.email IS 'Unique login email for staff member.';
COMMENT ON COLUMN staffs.password_hash IS 'Password hash for staff authentication.';
COMMENT ON COLUMN staffs.roles IS 'Assigned staff roles; must be unique and non-empty.';
COMMENT ON COLUMN staffs.is_active IS 'Whether staff account is allowed to authenticate.';
COMMENT ON COLUMN staffs.created_at IS 'Row creation timestamp.';
COMMENT ON COLUMN staffs.updated_at IS 'Row last update timestamp.';

-- +goose Down
DROP TRIGGER IF EXISTS staffs_set_updated_at ON staffs;
DROP INDEX IF EXISTS staffs_roles_idx;
DROP INDEX IF EXISTS staffs_email_uidx;
DROP TABLE IF EXISTS staffs;
DROP FUNCTION IF EXISTS staffs_set_updated_at();
DROP FUNCTION IF EXISTS staff_roles_is_valid(staff_role[]);
DROP TYPE IF EXISTS staff_role;
