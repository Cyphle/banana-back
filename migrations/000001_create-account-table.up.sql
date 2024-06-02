CREATE TABLE IF NOT EXISTS accounts(
    id          BIGSERIAL,
    name        VARCHAR(255),
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    deleted_at  timestamptz
    PRIMARY KEY (id)
)