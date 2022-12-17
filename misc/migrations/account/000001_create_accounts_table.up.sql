CREATE TABLE IF NOT EXISTS accounts (
    norek varchar(36) PRIMARY KEY,
    saldo float(32) NOT NULL DEFAULT 0,
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

INSERT INTO accounts (norek, saldo) VALUES ('r001', 10), ('r002', 5), ('r003', 2);