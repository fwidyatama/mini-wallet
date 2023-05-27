CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE wallets(
    id uuid DEFAULT uuid_generate_v4(),
    owned_by uuid NOT NULL unique ,
    status int default 1,
    balance bigint default 0,
    enabled_at timestamp DEFAULT 'now()',
    disabled_at timestamp,
    PRIMARY KEY (id)
)