CREATE TABLE IF NOT EXISTS profiles(
    id          SERIAL,
    username    VARCHAR(255),
    email       VARCHAR(255),
    first_name  VARCHAR(255),
    last_name   VARCHAR(255),
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    deleted_at  timestamptz,
    PRIMARY KEY (id)
)