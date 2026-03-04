-- +goose Up
CREATE TYPE domain_contact_role AS ENUM ('registrant', 'admin', 'tech', 'billing');

CREATE TABLE IF NOT EXISTS domains_contacts
(
    id         BIGSERIAL PRIMARY KEY,
    domain_id  BIGINT                                        NOT NULL,
    contact_id BIGINT                                        NOT NULL,
    role       domain_contact_role DEFAULT 'registrant'      NOT NULL,
    created_at TIMESTAMPTZ         DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_domain_id FOREIGN KEY (domain_id) REFERENCES domains (id),
    CONSTRAINT fk_contact_id FOREIGN KEY (contact_id) REFERENCES contacts (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS domain_contact_role_uidx
    ON domains_contacts (domain_id, contact_id, role);
CREATE UNIQUE INDEX IF NOT EXISTS domain_role_registrant_uidx
    ON domains_contacts (domain_id, role) WHERE role = 'registrant';

-- +goose Down
DROP INDEX IF EXISTS domain_role_registrant_uidx;
DROP INDEX IF EXISTS domain_contact_role_uidx;
DROP TABLE IF EXISTS domains_contacts;
DROP TYPE IF EXISTS domain_contact_role;
