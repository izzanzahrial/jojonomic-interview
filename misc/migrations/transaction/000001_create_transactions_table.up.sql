CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    norek varchar(36) NOT NULL,
    types varchar(36) NOT NULL,
    gram float(32) NOT NULL,
    topup_price integer NOT NULL,
    buyback_price integer NOT NULL,
    saldo float(32) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);