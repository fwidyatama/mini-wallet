CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE histories
(
    id uuid DEFAULT uuid_generate_v4(),
    wallet_id uuid NOT NULL,
    status VARCHAR,
    transaction_by uuid,
    type VARCHAR NOT NULL,
    amount FLOAT NOT NULL DEFAULT 0,
    reference_id uuid NOT NULL unique ,
    transaction_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_wallet_history_id FOREIGN KEY (wallet_id) REFERENCES "wallets" (id)
);