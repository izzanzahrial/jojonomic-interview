CREATE TABLE IF NOT EXISTS prices (
    admin_id varchar(36) PRIMARY KEY,
    topup_price integer NOT NULL,
    buyback_price integer NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);