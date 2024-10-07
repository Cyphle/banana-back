CREATE TABLE IF NOT EXISTS accounts(
    id                  BIGSERIAL,
    name                VARCHAR(255),
    type                VARCHAR(10),
    starting_balance    DECIMAL,
    profile_id          INTEGER,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    deleted_at          timestamptz,
    PRIMARY KEY (id),
    CONSTRAINT fk_accounts_profile FOREIGN KEY(profile_id) REFERENCES profiles(id)
)