-- +goose Up
CREATE TABLE IF NOT EXISTS registrars
(
    id                  BIGSERIAL PRIMARY KEY,
    username            TEXT                                  NOT NULL,
    password_hash       TEXT                                  NOT NULL,
    email               TEXT                                  NOT NULL,
    is_active           BOOLEAN     DEFAULT TRUE              NOT NULL,
    max_active_sessions INT         DEFAULT 1                 NOT NULL,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT max_active_sessions_check CHECK (max_active_sessions > 0),
    CONSTRAINT username_check CHECK (
        char_length(username) BETWEEN 3 AND 255 AND lower(username) ~* '^[a-z0-9]+$'
        ),
    CONSTRAINT email_check CHECK (
        char_length(email) >= 5 AND lower(email) ~* '^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$'
        )
);

CREATE UNIQUE INDEX IF NOT EXISTS registrars_username_uidx ON registrars (lower(username));
CREATE UNIQUE INDEX IF NOT EXISTS registrars_email_uidx ON registrars (lower(email));

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION registrars_set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

CREATE TRIGGER registrars_updated_at_trigger
    BEFORE UPDATE
    ON registrars
    FOR EACH ROW
EXECUTE PROCEDURE registrars_set_updated_at();

COMMENT ON TABLE registrars IS 'EPP registrars (clients) allowed to authenticate and manage registry objects.';
COMMENT ON COLUMN registrars.id IS 'Internal surrogate key.';
COMMENT ON COLUMN registrars.username IS 'Registrar login name used as EPP client identifier.';
COMMENT ON COLUMN registrars.password_hash IS 'Password hash for registrar authentication.';
COMMENT ON COLUMN registrars.email IS 'Registrar operational email.';
COMMENT ON COLUMN registrars.is_active IS 'Indicates whether registrar account can authenticate.';
COMMENT ON COLUMN registrars.max_active_sessions IS 'Maximum simultaneously active EPP sessions for registrar.';
COMMENT ON COLUMN registrars.created_at IS 'Row creation timestamp.';
COMMENT ON COLUMN registrars.updated_at IS 'Row last update timestamp.';

-- +goose Down
DROP TRIGGER IF EXISTS registrars_updated_at_trigger ON registrars;
DROP FUNCTION IF EXISTS registrars_set_updated_at();
DROP INDEX IF EXISTS registrars_email_uidx;
DROP INDEX IF EXISTS registrars_username_uidx;
DROP TABLE IF EXISTS registrars;
